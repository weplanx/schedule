package task

import (
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron/v3"
	"schedule-microservice/common"
	"time"
)

type Task struct {
	runtime map[string]*cron.Cron
	options map[string]*common.TaskOption
	entries map[string]map[string]cron.EntryID
}

func Create() *Task {
	task := new(Task)
	task.runtime = make(map[string]*cron.Cron)
	task.options = make(map[string]*common.TaskOption)
	task.entries = make(map[string]map[string]cron.EntryID)
	return task
}

func (c *Task) close(identity string) {
	if c.runtime[identity] == nil {
		return
	}
	if c.options[identity] != nil {
		for _, entry := range c.runtime[identity].Entries() {
			c.runtime[identity].Remove(entry.ID)
		}
	}
	c.runtime[identity].Stop()
}

func (c *Task) Get(identity string) *common.TaskOption {
	if c.options[identity] == nil || c.runtime[identity] == nil || c.entries[identity] == nil {
		return nil
	}
	options := c.options[identity]
	for key, entryId := range c.entries[identity] {
		entry := c.runtime[identity].Entry(entryId)
		options.Entries[key].NextDate = entry.Next
		options.Entries[key].LastDate = entry.Prev
	}
	return c.options[identity]
}

func (c *Task) Put(option common.TaskOption) (err error) {
	identity := option.Identity
	timezone, err := time.LoadLocation(option.TimeZone)
	if err != nil {
		return
	}
	c.close(identity)
	c.options[identity] = &option
	c.runtime[identity] = cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	c.entries[identity] = make(map[string]cron.EntryID)
	for key := range option.Entries {
		go c.webhook(identity, key)
	}
	if option.Start {
		c.runtime[identity].Start()
	}
	return
}

func (c *Task) webhook(identity string, key string) {
	var err error
	option := c.options[identity].Entries[key]
	c.entries[identity][key], err = c.runtime[identity].AddFunc(option.CronTime, func() {
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
			var message []string
			for _, value := range errs {
				message = append(message, value.Error())
			}
			//common.Record <- common.RecordError{
			//	Type:     "error",
			//	Identity: identity,
			//	Key:      key,
			//	Url:      option.Url,
			//	Header:   option.Headers,
			//	Body:     option.Body,
			//	Message:  message,
			//	Time:     time.Now().Unix(),
			//}
		} else {
			var response interface{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				println(err.Error())
			} else {
				//common.Record <- common.RecordSuccess{
				//	Type:     "success",
				//	Identity: identity,
				//	Key:      key,
				//	Url:      option.Url,
				//	Header:   option.Headers,
				//	Body:     option.Body,
				//	Response: response,
				//	Time:     time.Now().Unix(),
				//}
			}
		}
	})
	if err != nil {
		return
	}
}

func (c *Task) Delete(identity string) (err error) {
	c.close(identity)
	delete(c.runtime, identity)
	delete(c.options, identity)
	delete(c.entries, identity)
	return
}

func (c *Task) Running(identity string, running bool) (err error) {
	if c.runtime[identity] == nil {
		err = errors.New("this identity does not exists")
		return
	}
	if running {
		c.runtime[identity].Start()
	} else {
		c.runtime[identity].Stop()
	}
	return
}
