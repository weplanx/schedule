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
	// 监听同步
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

	// 同步定时发布
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

	// 订阅事件状态
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

type Sync struct {
	Key  string
	Time time.Time
}

func (x *App) PubSync(key string) (err error) {
	subject := fmt.Sprintf(`%s.sync`, x.Values.Namespace)
	var b []byte
	if b, err = msgpack.Marshal(Sync{
		Key:  key,
		Time: time.Now().Add(time.Second * 5),
	}); err != nil {
		return
	}
	if err = x.Nats.Publish(subject, b); err != nil {
		return
	}
	return
}
