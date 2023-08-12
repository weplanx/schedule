package common

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Inject struct {
	Values   *Values
	Log      *zap.Logger
	Nats     *nats.Conn
	KeyValue nats.KeyValue
}

type Values struct {
	Id string `env:"ID,required"`

	Nats struct {
		Hosts []string `env:"HOSTS,required" envSeparator:","`
		Nkey  string   `env:"NKEY,required"`
	} `envPrefix:"NATS_"`
}
