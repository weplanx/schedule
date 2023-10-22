package common

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

type Inject struct {
	V        *Values
	Log      *zap.Logger
	Nats     *nats.Conn
	KeyValue nats.KeyValue
}

type Values struct {
	Node string `env:"NODE,required"`
	Nats struct {
		Hosts []string `env:"HOSTS,required" envSeparator:","`
		Nkey  string   `env:"NKEY,required"`
	} `envPrefix:"NATS_"`
}

type Job struct {
	Key    string      `msgpack:"key"`
	Index  int         `msgpack:"index"`
	Mode   string      `msgpack:"mode"`
	Option interface{} `msgpack:"option"`
}

type HttpOption struct {
	Method  string                 `json:"method" msgpack:"method"`
	Url     string                 `json:"url" msgpack:"url"`
	Headers map[string]string      `json:"headers" msgpack:"headers"`
	Body    map[string]interface{} `json:"body" msgpack:"body"`
}

type ScheduleOption struct {
	Status bool          `json:"status" msgpack:"status" `
	Jobs   []ScheduleJob `json:"jobs" msgpack:"jobs" `
}

type ScheduleJob struct {
	Mode          string      `json:"mode" msgpack:"mode"`
	Spec          string      `json:"spec" msgpack:"spec"`
	Option        interface{} `json:"option" msgpack:"option"`
	ScheduleState `json:"schedule_state" msgpack:"state"`
}

type ScheduleState struct {
	Next time.Time `json:"next" msgpack:"next"`
	Prev time.Time `json:"prev" msgpack:"prev"`
}

type ScheduleStatus struct {
	Key   string `msgpack:"key"`
	Value bool   `msgpack:"value"`
}
