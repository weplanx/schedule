package client

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/utiliy"
)

type Schedule struct {
	Namespace string
	Js        nats.JetStreamContext
	Store     nats.ObjectStore
}

func New(namespace string, js nats.JetStreamContext) (x *Schedule, err error) {
	x = &Schedule{
		Namespace: namespace,
		Js:        js,
	}
	if x.Store, err = js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: fmt.Sprintf(`%s_schedules`, x.Namespace),
	}); err != nil {
		return
	}
	return
}

// Get 获取调度信息
func (x *Schedule) Get(key string) {

}

// Set 设置调度
func (x *Schedule) Set(key string, jobs ...utiliy.Job) (err error) {
	var b []byte
	if b, err = msgpack.Marshal(jobs); err != nil {
		return
	}
	if _, err = x.Store.PutBytes(key, b); err != nil {
		return
	}
	return
}

// Remove 移除调度
func (x *Schedule) Remove(key string) (err error) {
	if err = x.Store.Delete(key); err != nil {
		return
	}
	return
}
