package types

import "time"

type EntryOption struct {
	CronTime string            `yaml:"cron_time"`
	Url      string            `yaml:"url"`
	Headers  map[string]string `yaml:"headers"`
	Body     interface{}       `yaml:"body"`
	NextDate time.Time         `yaml:"-"`
	LastDate time.Time         `yaml:"-"`
}
