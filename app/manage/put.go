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
	return
}

func (c *JobsManager) addTask(identity string, taskID string) {
	option := c.options.Get(identity).Entries[taskID]
	EntryID, err := c.runtime.Get(identity).AddFunc(option.CronTime, func() {
		actions.Fetch(types.FetchOption{
			Url:     option.Url,
			Headers: option.Headers,
			Body:    option.Body,
		})
	})
	if err != nil {
		return
	}
	c.entryIDSet.Set(identity, taskID, EntryID)
}
