package job

import (
	"errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"schedule-microservice/application/service/filelog"
	"schedule-microservice/application/service/job/utils"
	"schedule-microservice/application/service/schema"
	"schedule-microservice/application/service/transfer"
	"schedule-microservice/config/options"
)

type Job struct {
	Options    *utils.JobOptionMap
	Runtime    *utils.CronMap
	EntryIDSet *utils.EntryIDMap
	*Dependency
}

var (
	NotExists = errors.New("this identity does not exists")
)

type Dependency struct {
	fx.In

	Schema   *schema.Schema
	Filelog  *filelog.Filelog
	Transfer *transfer.Transfer
}

func New(dep *Dependency) (c *Job, err error) {
	c = new(Job)
	c.Dependency = dep
	c.Options = utils.NewJobOptionMap()
	c.Runtime = utils.NewCronMap()
	c.EntryIDSet = utils.NewEntryIDMap()
	var jobOptions []options.JobOption
	if jobOptions, err = c.Schema.Lists(); err != nil {
		return
	}
	for _, option := range jobOptions {
		if err = c.Put(option); err != nil {
			return
		}
	}
	return
}

func (c *Job) Get(identity string) (*options.JobOption, *cron.Cron, map[string]cron.EntryID, error) {
	if c.Options.Empty(identity) {
		return nil, nil, nil, NotExists
	}
	return c.Options.Get(identity), c.Runtime.Get(identity), c.EntryIDSet.Get(identity), nil
}

func (c *Job) termination(identity string) {
	if c.Options.Empty(identity) {
		return
	}
	runtime := c.Runtime.Get(identity)
	for _, entry := range runtime.Entries() {
		runtime.Remove(entry.ID)
	}
	runtime.Stop()
}
