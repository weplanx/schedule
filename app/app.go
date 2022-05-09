package app

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/utiliy"
	"go.uber.org/zap"
	"time"
)

// Run 启动服务
func (x *App) Run() (err error) {
	// 多节点工作同步
	if err = x.Sync(); err != nil {
		return
	}
	// 状态反馈
	if err = x.State(); err != nil {
		return
	}
	// 拉取配置设置定时
	var objects []*nats.ObjectInfo
	if objects, err = x.Store.List(); errors.Is(err, nats.ErrNoObjectsFound) {
		if errors.Is(err, nats.ErrNoObjectsFound) {
			objects = make([]*nats.ObjectInfo, 0)
		} else {
			return
		}
	}
	for _, o := range objects {
		key := o.Name
		var b []byte
		if b, err = x.Store.GetBytes(key); err != nil {
			return
		}
		var jobs []utiliy.Job
		if err = msgpack.Unmarshal(b, &jobs); err != nil {
			x.Log.Error("解码失败",
				zap.ByteString("data", b),
				zap.Error(err),
			)
			return
		}
		if err = x.Schedules.Set(key, jobs...); err != nil {
			x.Log.Error("定时发布设置失败",
				zap.String("key", key),
				zap.Any("jobs", jobs),
				zap.Error(err),
			)
			return
		}
		if err = x.PubSync(key); err != nil {
			return
		}
	}
	// 订阅配置
	var watch nats.ObjectWatcher
	if watch, err = x.Store.Watch(); err != nil {
		return
	}
	current := time.Now()
	for o := range watch.Updates() {
		if o == nil || o.ModTime.Unix() < current.Unix() {
			continue
		}
		key := o.Name
		if !o.Deleted {
			var b []byte
			if b, err = x.Store.GetBytes(key); err != nil {
				return
			}
			var jobs []utiliy.Job
			if err = msgpack.Unmarshal(b, &jobs); err != nil {
				x.Log.Error("解码失败",
					zap.ByteString("data", b),
					zap.Error(err),
				)
				return
			}
			if err = x.Schedules.Set(key, jobs...); err != nil {
				x.Log.Error("定时发布设置失败",
					zap.String("key", key),
					zap.Any("jobs", jobs),
					zap.Error(err),
				)
				return
			}
			if err = x.PubSync(key); err != nil {
				return
			}
		} else {
			x.Schedules.Remove(key)
		}
	}
	return
}

func (x *App) Sync() (err error) {
	subject := fmt.Sprintf(`%s.sync`, x.Values.Namespace)
	if _, err = x.Nats.Subscribe(subject, func(msg *nats.Msg) {
		var sync Sync
		if err := msgpack.Unmarshal(msg.Data, &sync); err != nil {
			return
		}
		x.Schedules.Stop(sync.Key)
		time.Sleep(sync.Time.Sub(time.Now()))
		x.Schedules.Start(sync.Key)
	}); err != nil {
		return
	}
	return
}

func (x *App) State() (err error) {
	name := fmt.Sprintf(`%s:state`, x.Values.Namespace)
	subject := fmt.Sprintf(`%s.state`, x.Values.Namespace)
	if _, err = x.Nats.QueueSubscribe(subject, name, func(msg *nats.Msg) {
		key := string(msg.Data)
		var values []utiliy.State
		for _, entry := range x.Schedules.State(key) {
			values = append(values, utiliy.State{
				Next: entry.Next,
				Prev: entry.Prev,
			})
		}
		b, _ := msgpack.Marshal(values)
		msg.Respond(b)
	}); err != nil {
		return
	}
	return
}
