package common

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Inject struct {
	Values *Values
	Log    *zap.Logger
	Nats   *nats.Conn
	Js     nats.JetStreamContext
	Store  nats.ObjectStore
}
