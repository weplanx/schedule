package typ

import "time"

type Job struct {
	Key    string      `msgpack:"key"`
	Index  int         `msgpack:"index"`
	Mode   string      `msgpack:"mode"`
	Option interface{} `msgpack:"option"`
}

type HttpOption struct {
	Url     string                 `msgpack:"url"`
	Headers map[string]string      `msgpack:"headers"`
	Body    map[string]interface{} `msgpack:"body"`
}

type ScheduleJob struct {
	Mode          string                 `msgpack:"mode"`
	Spec          string                 `msgpack:"spec"`
	Option        map[string]interface{} `msgpack:"option"`
	ScheduleState `msgpack:"state"`
}

type ScheduleState struct {
	Next time.Time `msgpack:"next"`
	Prev time.Time `msgpack:"prev"`
}

type ScheduleStatus struct {
	Key     string `msgpack:"key"`
	Running bool   `msgpack:"running"`
}
