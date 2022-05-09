package utiliy

import (
	"github.com/robfig/cron/v3"
	"log"
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
	if x.values[key] != nil {
		x.Remove(key)
	}
	x.values[key] = cron.New(cron.WithSeconds())
	for _, job := range jobs {
		if _, err = x.values[key].AddFunc(job.Rule, func() {
			log.Println(key, job.Rule, job.Mode, job.Spec)
		}); err != nil {
			return
		}
	}
	return
}

func (x *Schedule) Start(key string) {
	x.values[key].Start()
}

func (x *Schedule) Stop(key string) {
	x.values[key].Stop()
}

func (x *Schedule) State(key string) []cron.Entry {
	return x.values[key].Entries()
}

func (x *Schedule) Remove(key string) {
	if c, exists := x.values[key]; exists {
		c.Stop()
		delete(x.values, key)
	}
}
