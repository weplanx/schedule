package common

import (
	"github.com/nats-io/nats.go"
	"github.com/weplanx/collector/transfer"
	"go.uber.org/zap"
)

type Inject struct {
	V         *Values
	Log       *zap.Logger
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	Transfer  *transfer.Transfer
}

type Values struct {
	Namespace string `env:"NAMESPACE,required"`
	Nats      struct {
		Hosts []string `env:"HOSTS,required" envSeparator:","`
		Nkey  string   `env:"NKEY,required"`
	} `envPrefix:"NATS_"`
}
