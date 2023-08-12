package typ

import "time"

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
	Key     string `msgpack:"key"`
	Running bool   `msgpack:"running"`
}
