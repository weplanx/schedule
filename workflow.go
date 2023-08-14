package workflow

import (
	"github.com/nats-io/nats.go"
)

type Workflow struct {
	Namespace string
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
}

func New(options ...Option) *Workflow {
	x := new(Workflow)
	for _, v := range options {
		v(x)
	}
	return x
}

type Option func(x *Workflow)

func SetNamespace(v string) Option {
	return func(x *Workflow) {
		x.Namespace = v
	}
}

func SetNats(v *nats.Conn) Option {
	return func(x *Workflow) {
		x.Nats = v
	}
}

func SetJetStream(v nats.JetStreamContext) Option {
	return func(x *Workflow) {
		x.JetStream = v
	}
}
