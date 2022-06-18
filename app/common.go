package app

import (
	"fmt"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/common"
	"go.uber.org/zap"
	"time"
)

var Provides = wire.NewSet(New)

type App struct {
	*common.Inject
	values map[string]*cron.Cron
}

func New(i *common.Inject) *App {
	return &App{
		Inject: i,
		values: make(map[string]*cron.Cron),
	}
}

// Set 设置任务发布
func (x *App) Set(key string, jobs ...common.Job) (err error) {
	if x.values[key] != nil {
		x.Remove(key)
	}
	x.values[key] = cron.New(cron.WithSeconds())
	for n, job := range jobs {
		if _, err = x.values[key].AddFunc(job.Spec, func() {
			subject := fmt.Sprintf(`%s.schedules`, x.Values.Namespace)
			value := map[string]interface{}{
				"key":    key,
				"n":      n,
				"mode":   job.Mode,
				"option": job.Option,
			}
			b, _ := msgpack.Marshal(value)
			msgId := fmt.Sprintf(`%s-%d`, key, time.Now().Unix())
			if _, err := x.Js.Publish(subject, b, nats.MsgId(msgId)); err != nil {
				x.Log.Error("发布失败",
					zap.String("key", key),
					zap.Int("n", n),
					zap.Error(err),
				)
				return
			}
			x.Log.Info("任务发布成功",
				zap.String("key", key),
				zap.String("msgId", msgId),
				zap.Int("n", n),
				zap.String("mode", job.Mode),
				zap.Any("option", job.Option),
			)
		}); err != nil {
			return
		}
	}
	return
}

// Start 启动任务
func (x *App) Start(key string) {
	if c, ok := x.values[key]; ok {
		c.Start()
	}
}

// Stop 关闭任务
func (x *App) Stop(key string) {
	if c, ok := x.values[key]; ok {
		c.Stop()
	}
}

// State 任务状态
func (x *App) State(key string) []cron.Entry {
	if c, ok := x.values[key]; ok {
		return c.Entries()
	}
	return []cron.Entry{}
}

// Remove 移除任务
func (x *App) Remove(key string) {
	if c, ok := x.values[key]; ok {
		c.Stop()
		delete(x.values, key)
	}
}

type Sync struct {
	Key  string
	Time time.Time
}

// PubSync 发布同步工作
func (x *App) PubSync(key string) (err error) {
	subject := fmt.Sprintf(`%s.sync`, x.Values.Namespace)
	var b []byte
	if b, err = msgpack.Marshal(Sync{
		Key:  key,
		Time: time.Now().Add(time.Second * 3),
	}); err != nil {
		return
	}
	if err = x.Nats.Publish(subject, b); err != nil {
		return
	}
	return
}
