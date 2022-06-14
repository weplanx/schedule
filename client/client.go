package client

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/common"
	"time"
)

type Schedule struct {
	Namespace string
	Nats      *nats.Conn
	Js        nats.JetStreamContext
	Store     nats.ObjectStore
}

// New 创建客户端
func New(namespace string, nc *nats.Conn, js nats.JetStreamContext) (x *Schedule, err error) {
	x = &Schedule{
		Namespace: namespace,
		Nats:      nc,
		Js:        js,
	}
	if x.Store, err = js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: fmt.Sprintf(`%s_schedules`, x.Namespace),
	}); err != nil {
		return
	}
	return
}

// List 列出配置标识
func (x *Schedule) List() (keys []string, err error) {
	var infos []*nats.ObjectInfo
	if infos, err = x.Store.List(); err != nil {
		if errors.Is(err, nats.ErrNoObjectsFound) {
			return []string{}, nil
		}
		return
	}
	keys = make([]string, len(infos))
	for i, x := range infos {
		keys[i] = x.Name
	}
	return
}

// Get 获取调度信息
func (x *Schedule) Get(key string) (jobs []common.Job, err error) {
	var b []byte
	if b, err = x.Store.GetBytes(key); err != nil {
		return
	}
	if err = msgpack.Unmarshal(b, &jobs); err != nil {
		return
	}
	subject := fmt.Sprintf(`%s.state`, x.Namespace)
	var msg *nats.Msg
	if msg, err = x.Nats.Request(subject, []byte(key), time.Second*3); err != nil {
		return
	}
	var values []common.State
	if err = msgpack.Unmarshal(msg.Data, &values); err != nil {
		return
	}
	for k, v := range values {
		jobs[k].State = v
	}
	return
}

// Set 设置调度
func (x *Schedule) Set(key string, jobs ...common.Job) (err error) {
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
