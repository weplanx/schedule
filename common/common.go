package common

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

type Inject struct {
	Values *Values
	Log    *zap.Logger
	Nats   *nats.Conn
	Js     nats.JetStreamContext
	Store  nats.ObjectStore
}

type Values struct {
	Namespace string `env:"NAMESPACE,required"`
	Nats      `envPrefix:"NATS_"`
}

type Nats struct {
	Hosts []string `env:"HOSTS,required" envSeparator:","`
	Nkey  string   `env:"NKEY,required"`
}

type Job struct {
	Mode   string `msgpack:"mode"`
	Spec   string `msgpack:"spec"`
	Option Option `msgpack:"option"`
	State  `msgpack:"state"`
}

type Option interface{}

type HttpOption struct {
	Url     string                 `msgpack:"url"`
	Headers map[string]string      `msgpack:"headers"`
	Body    map[string]interface{} `msgpack:"body"`
}

func HttpJob(spec string, option HttpOption) Job {
	return Job{
		Mode:   "HTTP",
		Spec:   spec,
		Option: option,
	}
}

type State struct {
	Next time.Time `msgpack:"next"`
	Prev time.Time `msgpack:"prev"`
}

type Status struct {
	Key     string `msgpack:"key"`
	Running bool   `msgpack:"running"`
}
