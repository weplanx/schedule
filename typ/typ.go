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

type ScheduleOption struct {
	Status bool          `msgpack:"status"`
	Jobs   []ScheduleJob `msgpack:"jobs"`
}

type ScheduleJob struct {
	Mode          string      `msgpack:"mode"`
	Spec          string      `msgpack:"spec"`
	Option        interface{} `msgpack:"option"`
	ScheduleState `msgpack:"state"`
}

type ScheduleState struct {
	Next time.Time `msgpack:"next"`
	Prev time.Time `msgpack:"prev"`
}

type ScheduleStatus struct {
	Key   string `msgpack:"key"`
	Value bool   `msgpack:"value"`
}

type ExcelMetadata struct {
	Name  string   `msgpack:"name"`
	Parts []string `msgpack:"parts"`
}
type ExcelSheets map[string][][]interface{}
