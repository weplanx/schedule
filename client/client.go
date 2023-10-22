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
	Node      string
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	KeyValue  nats.KeyValue
}

func New(node string, nc *nats.Conn) (x *Client, err error) {
	x = &Client{Node: node, Nats: nc}
	if x.JetStream, err = nc.JetStream(
		nats.PublishAsyncMaxPending(256),
	); err != nil {
		return
	}
	bucket := fmt.Sprintf(`schedules_%s`, x.Node)
	if x.KeyValue, err = x.JetStream.KeyValue(bucket); err != nil {
		return
	}
	return
}

func (x *Client) Ping() (result bool, err error) {
	var msg *nats.Msg
	if msg, err = x.Nats.Request("schedules", []byte(x.Node), time.Second*5); err != nil {
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
	subj := fmt.Sprintf(`schedules.%s`, x.Node)
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
