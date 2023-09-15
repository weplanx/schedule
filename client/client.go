package client

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/common"
	"time"
)

type Client struct {
	Namespace string
	Node      string
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	KeyValue  nats.KeyValue
}

type Option func(x *Client)

func SetNamespace(v string) Option {
	return func(x *Client) {
		x.Namespace = v
	}
}

func SetNode(v string) Option {
	return func(x *Client) {
		x.Node = v
	}
}

func SetNats(v *nats.Conn) Option {
	return func(x *Client) {
		x.Nats = v
	}
}

func SetJetStream(v nats.JetStreamContext) Option {
	return func(x *Client) {
		x.JetStream = v
	}
}

func New(options ...Option) (x *Client, err error) {
	x = new(Client)
	for _, apply := range options {
		apply(x)
	}
	bucket := fmt.Sprintf(`%s_schedules_%s`, x.Namespace, x.Node)
	if x.KeyValue, err = x.JetStream.KeyValue(bucket); err != nil {
		return
	}
	return
}

func (x *Client) Ping() (result bool, err error) {
	subj := fmt.Sprintf(`%s.schedules`, x.Namespace)
	var msg *nats.Msg
	if msg, err = x.Nats.Request(subj, []byte(x.Node), time.Second*5); err != nil {
		return
	}
	result = string(msg.Data) == "ok"
	return
}

func (x *Client) Lists() (keys []string, err error) {
	if keys, err = x.KeyValue.Keys(); err != nil {
		if errors.Is(err, nats.ErrNoKeysFound) {
			return []string{}, nil
		}
		return
	}
	return
}

func (x *Client) Get(key string) (option common.ScheduleOption, err error) {
	var entry nats.KeyValueEntry
	if entry, err = x.KeyValue.Get(key); err != nil {
		return
	}
	if err = msgpack.Unmarshal(entry.Value(), &option); err != nil {
		return
	}
	var msg *nats.Msg
	subj := fmt.Sprintf(`%s.schedules.%s`, x.Namespace, x.Node)
	if msg, err = x.Nats.Request(subj, []byte(key), time.Second*3); err != nil {
		return
	}
	var states []common.ScheduleState
	if err = msgpack.Unmarshal(msg.Data, &states); err != nil {
		return
	}
	for k, v := range states {
		option.Jobs[k].ScheduleState = v
	}
	return
}

func (x *Client) Set(key string, option common.ScheduleOption) (err error) {
	var b []byte
	if b, err = msgpack.Marshal(option); err != nil {
		return
	}
	if _, err = x.KeyValue.Put(key, b); err != nil {
		return
	}
	return
}

func (x *Client) Status(key string, value bool) (err error) {
	var entry nats.KeyValueEntry
	if entry, err = x.KeyValue.Get(key); err != nil {
		return
	}
	var option common.ScheduleOption
	if err = msgpack.Unmarshal(entry.Value(), &option); err != nil {
		return
	}
	option.Status = value
	return x.Set(key, option)
}

func (x *Client) Remove(key string) (err error) {
	return x.KeyValue.Delete(key)
}
