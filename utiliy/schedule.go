package utiliy

import (
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	values map[string]*cron.Cron
}

func NewSchedule() *Schedule {
	return &Schedule{
		values: make(map[string]*cron.Cron),
	}
}

func (x *Schedule) Set(key string, jobs ...Job) (err error) {
	x.values[key] = cron.New(cron.WithSeconds())
	for _, job := range jobs {
		if _, err = x.values[key].AddFunc(job.Rule, func() {
		}); err != nil {
			return
		}
	}
	return
}

func (x *Schedule) Start(k string) {
	x.values[k].Start()
}

func (x *Schedule) State(k string) []cron.Entry {
	return x.values[k].Entries()
}

func (x *Schedule) Remove(k string) {
	if c, exists := x.values[k]; exists {
		c.Stop()
		delete(x.values, k)
	}
}
