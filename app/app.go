package app

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/common"
	"go.uber.org/zap"
	"time"
)

func (x *App) Run() (err error) {
	name := fmt.Sprintf(`%s:schedules`, x.Values.Namespace)
	subject := fmt.Sprintf(`%s.schedules`, x.Values.Namespace)
	if _, err = x.Js.AddStream(&nats.StreamConfig{
		Name:        name,
		Subjects:    []string{subject},
		Description: "定时调度发布",
	}); err != nil {
		return
	}
	if err = x.SubSync(); err != nil {
		return
	}
	if err = x.SubState(); err != nil {
		return
	}
	if err = x.SubStatus(); err != nil {
		return
	}
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
		var jobs []common.Job
		if err = msgpack.Unmarshal(b, &jobs); err != nil {
			x.Log.Error("解码失败",
				zap.ByteString("data", b),
				zap.Error(err),
			)
			return
		}
		if err = x.Set(key, jobs...); err != nil {
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
		x.Log.Info("已设置发布",
			zap.String("key", key),
			zap.Any("jobs", jobs),
		)
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
			var jobs []common.Job
			if err = msgpack.Unmarshal(b, &jobs); err != nil {
				x.Log.Error("解码失败",
					zap.ByteString("data", b),
					zap.Error(err),
				)
				return
			}
			if err = x.Set(key, jobs...); err != nil {
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
			x.Log.Info("定时发布成功",
				zap.String("key", key),
				zap.Any("jobs", jobs),
			)
		} else {
			x.Remove(key)
			x.Log.Info("定时移除成功",
				zap.String("key", key),
			)
		}
	}
	return
}

func (x *App) SubSync() (err error) {
	subject := fmt.Sprintf(`%s.sync`, x.Values.Namespace)
	if _, err = x.Nats.Subscribe(subject, func(msg *nats.Msg) {
		var sync Sync
		if err := msgpack.Unmarshal(msg.Data, &sync); err != nil {
			return
		}
		x.Stop(sync.Key)
		time.Sleep(sync.Time.Sub(time.Now()))
		x.Start(sync.Key)
		x.Log.Info("节点时间已同步")
	}); err != nil {
		return
	}
	return
}

func (x *App) SubState() (err error) {
	name := fmt.Sprintf(`%s:state`, x.Values.Namespace)
	subject := fmt.Sprintf(`%s.state`, x.Values.Namespace)
	if _, err = x.Nats.QueueSubscribe(subject, name, func(msg *nats.Msg) {
		key := string(msg.Data)
		var values []common.State
		for _, entry := range x.State(key) {
			values = append(values, common.State{
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

func (x *App) SubStatus() (err error) {
	subject := fmt.Sprintf(`%s.status`, x.Values.Namespace)
	if _, err = x.Nats.Subscribe(subject, func(msg *nats.Msg) {
		var data common.Status
		if err := msgpack.Unmarshal(msg.Data, &data); err != nil {
			return
		}
		if data.Running {
			x.Start(data.Key)
			x.Log.Debug("任务状态已启动",
				zap.String("key", data.Key),
			)
		} else {
			x.Stop(data.Key)
			x.Log.Debug("任务状态已停止",
				zap.String("key", data.Key),
			)
		}
		msg.Respond([]byte("ok"))
	}); err != nil {
		return
	}
	return
}
