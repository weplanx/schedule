package manage

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"schedule-microservice/app/actions"
	"schedule-microservice/app/types"
	"time"
)

func (c *JobsManager) Put(option types.JobOption) (err error) {
	identity := option.Identity
	timezone, err := time.LoadLocation(option.TimeZone)
	if err != nil {
		return
	}
	c.termination(identity)
	c.options.Set(identity, &option)
	job := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	c.runtime.Set(identity, job)
	for taskID := range option.Entries {
		go c.addTask(identity, taskID)
	}
	if option.Start {
		job.Start()
	}
	return c.schema.Update(option)
}

func (c *JobsManager) addTask(identity string, taskID string) {
	option := c.options.Get(identity).Entries[taskID]
	EntryID, err := c.runtime.Get(identity).AddFunc(option.CronTime, func() {
		body, errs := actions.Fetch(types.FetchOption{
			Url:     option.Url,
			Headers: option.Headers,
			Body:    option.Body,
		})
		var message map[string]interface{}
		var bodyRecord interface{}
		bodyRaw, ok := option.Body.(string)
		if ok && json.Valid([]byte(bodyRaw)) {
			json.Unmarshal([]byte(bodyRaw), &bodyRecord)
		} else {
			bodyRecord = option.Body
		}
		if len(errs) != 0 {
			msg := make([]string, len(errs))
			for index, value := range errs {
				msg[index] = value.Error()
			}
			message = map[string]interface{}{
				"Identity": identity,
				"Task":     taskID,
				"Url":      option.Url,
				"Header":   option.Headers,
				"Body":     bodyRecord,
				"Msg":      msg,
				"Time":     time.Now().Unix(),
			}
		} else {
			var responseRecord interface{}
			if json.Valid(body) {
				json.Unmarshal(body, &responseRecord)
			}
			message = map[string]interface{}{
				"Identity": identity,
				"Task":     taskID,
				"Url":      option.Url,
				"Header":   option.Headers,
				"Body":     bodyRecord,
				"Response": responseRecord,
				"Time":     time.Now().Unix(),
			}
		}
		c.logging.Push(&types.LoggingPush{
			Identity: identity,
			HasError: len(errs) != 0,
			Message:  message,
		})
	})
	if err != nil {
		return
	}
	c.entryIDSet.Set(identity, taskID, EntryID)
}
