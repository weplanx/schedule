package job

import (
	"errors"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"os"
	"schedule-microservice/common"
	"time"
)

type Job struct {
	options *syncMapJobOption
	runtime *syncMapCron
	entries *syncMapEntryID
}

func Create() *Job {
	var err error
	c := new(Job)
	c.options = newSyncMapJobOption()
	c.runtime = newSyncMapCron()
	c.entries = newSyncMapEntryID()
	var configs []common.JobOption
	configs, err = common.ListConfig()
	if err != nil {
		logrus.Fatalln(err)
	}
	for _, opt := range configs {
		err = c.Put(opt)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
	return c
}

func (c *Job) termination(identity string) {
	if c.runtime.Empty(identity) || c.options.Empty(identity) {
		return
	}
	runtime := c.runtime.Get(identity)
	for _, entry := range runtime.Entries() {
		runtime.Remove(entry.ID)
	}
	runtime.Stop()
}

func (c *Job) Get(identity string) *common.JobOption {
	if c.options.Empty(identity) || c.runtime.Empty(identity) || c.entries.Empty(identity) {
		return nil
	}
	options := c.options.Map[identity]
	for key, entryId := range c.entries.GetJobEntryID(identity) {
		entry := c.runtime.Map[identity].Entry(entryId)
		options.Entries[key].NextDate = entry.Next
		options.Entries[key].LastDate = entry.Prev
	}
	return c.options.Map[identity]
}

func (c *Job) All() []string {
	var keys []string
	for key := range c.options.Map {
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
	c.options.Set(identity, &option)
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
	var err error
	option := c.options.Get(identity).Entries[task]
	EntryID, err := c.runtime.Get(identity).AddFunc(option.CronTime, func() {
		logger := logrus.New()
		var file *os.File
		if common.OpenStorage() {
			file, err = common.LogFile(identity)
			if err != nil {
				return
			}
			logger.SetOutput(file)
		}
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
			var msg []string
			for _, value := range errs {
				msg = append(msg, value.Error())
			}
			message := map[string]interface{}{
				"Identity": identity,
				"Task":     task,
				"Url":      option.Url,
				"Header":   option.Headers,
				"Body":     option.Body,
				"Msg":      msg,
				"Time":     time.Now().Unix(),
			}
			logger.Error(message)
			common.PushLogger(message)
		} else {
			message := map[string]interface{}{
				"Identity": identity,
				"Task":     task,
				"Url":      option.Url,
				"Header":   option.Headers,
				"Body":     option.Body,
				"Response": string(body),
				"Time":     time.Now().Unix(),
			}
			logger.Info(message)
			common.PushLogger(message)
		}
	})
	c.entries.Set(identity, task, EntryID)
	if err != nil {
		return
	}
}

func (c *Job) Delete(identity string) (err error) {
	c.termination(identity)
	c.options.Clear(identity)
	c.runtime.Clear(identity)
	c.entries.Clear(identity)
	return common.RemoveConfig(identity)
}

func (c *Job) Running(identity string, running bool) (err error) {
	if c.runtime.Empty(identity) || c.options.Empty(identity) {
		err = errors.New("this identity does not exists")
		return
	}
	c.options.Map[identity].Start = running
	if running {
		c.runtime.Map[identity].Start()
	} else {
		c.runtime.Map[identity].Stop()
	}
	return common.SaveConfig(c.options.Map[identity])
}
