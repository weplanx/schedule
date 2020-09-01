package manage

import (
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
				"Body":     option.Body,
				"Msg":      msg,
				"Time":     time.Now().Unix(),
			}
		} else {
			message = map[string]interface{}{
				"Identity": identity,
				"Task":     taskID,
				"Url":      option.Url,
				"Header":   option.Headers,
				"Body":     option.Body,
				"Response": string(body),
				"Time":     time.Now().Unix(),
			}
		}
		actions.Logging(c.logging, &types.LoggingPush{
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
