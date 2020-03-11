package common

import "time"

type (
	TaskOption struct {
		Identity string
		TimeZone string
		Start    bool
		Entries  map[string]*EntryOption
	}
	EntryOption struct {
		CronTime string
		Url      string
		Headers  map[string]string
		Body     interface{}
		NextDate time.Time
		LastDate time.Time
	}
)
