package job

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"log"
	"schedule-microservice/common"
	"time"
)

type Job struct {
	options map[string]*common.JobOption
	runtime *syncMapCron
	entries *syncMapEntryID
}

func Create() *Job {
	var err error
	c := new(Job)
	c.options = make(map[string]*common.JobOption)
	c.runtime = newSyncMapCron()
	c.entries = newSyncMapEntryID()
	var configs []common.JobOption
	configs, err = common.ListConfig()
	if err != nil {
		log.Fatalln(err)
	}
	for _, opt := range configs {
		err = c.Put(opt)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return c
}

func (c *Job) termination(identity string) {
	if c.runtime.Empty(identity) || c.options[identity] == nil {
		return
	}
	runtime := c.runtime.Get(identity)
	for _, entry := range runtime.Entries() {
		runtime.Remove(entry.ID)
	}
	runtime.Stop()
}

func (c *Job) Get(identity string) *common.JobOption {
	if c.options[identity] == nil || c.runtime.Empty(identity) || c.entries.Empty(identity) {
		return nil
	}
	options := c.options[identity]
	for key, entryId := range c.entries.GetJobEntryID(identity) {
		entry := c.runtime.Map[identity].Entry(entryId)
		options.Entries[key].NextDate = entry.Next
		options.Entries[key].LastDate = entry.Prev
	}
	return c.options[identity]
}

func (c *Job) All() []string {
	var keys []string
	for key := range c.options {
		keys = append(keys, key)
	}
	return keys
}

func (c *Job) Put(option common.JobOption) (err error) {
	identity := option.Identity
	timezone, err := time.LoadLocation(option.TimeZone)
	if err != nil {
		return
	}
	c.termination(identity)
	c.options[identity] = &option
	runtime := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	c.runtime.Set(identity, runtime)
	for key := range option.Entries {
		go c.addTask(identity, key)
	}
	if option.Start {
		runtime.Start()
	}
	return common.SaveConfig(&option)
}

func (c *Job) addTask(identity string, task string) {
	option := c.options[identity].Entries[task]
	EntryID, err := c.runtime.Get(identity).AddFunc(option.CronTime, func() {
		agent := gorequest.New().Post(option.Url)
		if option.Headers != nil {
			for key, value := range option.Headers {
				agent.Set(key, value)
			}
		}
		if option.Body != nil {
			agent.Send(option.Body)
		}
		_, body, errs := agent.EndBytes()
		if errs != nil {
			logrus.Info(errs)
			var message []string
			for _, value := range errs {
				message = append(message, value.Error())
			}
		} else {
			var response interface{}
			err := jsoniter.Unmarshal(body, &response)
			if err != nil {
				logrus.Error(err)
			} else {
			}
		}
	})
	c.entries.Set(identity, task, EntryID)
	if err != nil {
		return
	}
}

func (c *Job) Delete(identity string) (err error) {
	c.termination(identity)
	delete(c.options, identity)
	c.entries.Clear(identity)
	c.entries.Clear(identity)
	return common.RemoveConfig(identity)
}

func (c *Job) Running(identity string, running bool) (err error) {
	if c.runtime.Empty(identity) || c.options[identity] == nil {
		err = errors.New("this identity does not exists")
		return
	}
	c.options[identity].Start = running
	if running {
		c.runtime.Map[identity].Start()
	} else {
		c.runtime.Map[identity].Stop()
	}
	return common.SaveConfig(c.options[identity])
}
