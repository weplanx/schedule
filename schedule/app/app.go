package app

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/workflow/schedule/common"
	"github.com/weplanx/workflow/typ"
	"go.uber.org/zap"
	"sync"
	"time"
)

type App struct {
	*common.Inject
	M sync.Map
}

func Initialize(i *common.Inject) *App {
	return &App{
		Inject: i,
		M:      sync.Map{},
	}
}

func (x *App) Run() (err error) {
	var keys []string
	if keys, err = x.KeyValue.Keys(); err != nil {
		if errors.Is(err, nats.ErrNoKeysFound) {
			keys = []string{}
		} else {
			return
		}
	}
	for _, key := range keys {
		var entry nats.KeyValueEntry
		if entry, err = x.KeyValue.Get(key); err != nil {
			return
		}
		var option typ.ScheduleOption
		if err = msgpack.Unmarshal(entry.Value(), &option); err != nil {
			return
		}
		if err = x.Set(key, option); err != nil {
			return
		}
	}

	now := time.Now()
	var watch nats.KeyWatcher
	if watch, err = x.KeyValue.WatchAll(); err != nil {
		return
	}
	go func() {
		for update := range watch.Updates() {
			if update == nil || update.Created().Unix() < now.Unix() {
				continue
			}
			key := update.Key()
			switch update.Operation().String() {
			case "KeyValuePutOp":
				var option typ.ScheduleOption
				if err = msgpack.Unmarshal(update.Value(), &option); err != nil {
					return
				}
				if err = x.Set(key, option); err != nil {
					x.Log.Error("KeyValuePutOp Set:fail",
						zap.String("key", key),
						zap.Any("option", option),
						zap.Error(err),
					)
					return
				}
				x.Log.Debug("KeyValuePutOp SetCron:ok",
					zap.String("key", key),
					zap.Any("option", option),
				)
				break
			case "KeyValueDeleteOp":
				x.Remove(key)
				x.Log.Debug("KeyValueDeleteOp Remove:ok",
					zap.String("key", key),
				)
				break
			}
		}
	}()
	return
}

func (x *App) Get(key string) (*cron.Cron, bool) {
	v, ok := x.M.Load(key)
	return v.(*cron.Cron), ok
}

func (x *App) Set(key string, option typ.ScheduleOption) (err error) {
	x.Remove(key)
	c := cron.New(cron.WithSeconds())
	for i, job := range option.Jobs {
		if err = x.SetJob(key, c, i, job); err != nil {
			return
		}
	}
	if option.Status {
		c.Start()
	}
	x.M.Store(key, c)
	return
}

func (x *App) SetJob(key string, c *cron.Cron, index int, job typ.ScheduleJob) (err error) {
	if _, err = c.AddFunc(job.Spec, func() {
		subj := fmt.Sprintf(`%s.jobs.%s`, x.V.Namespace, x.V.Id)
		var b []byte
		if b, err = msgpack.Marshal(typ.Job{
			Key:    key,
			Index:  index,
			Mode:   job.Mode,
			Option: job.Option,
		}); err != nil {
			return
		}
		if err = x.Nats.Publish(subj, b); err != nil {
			x.Log.Error("Publish:fail",
				zap.String("key", key),
				zap.Int("index", index),
				zap.Error(err),
			)
			return
		}
		x.Log.Debug("Publish:ok",
			zap.String("key", key),
			zap.Int("index", index),
			zap.String("mode", job.Mode),
			zap.Any("option", job.Option),
		)
	}); err != nil {
		return
	}
	return
}

func (x *App) Remove(key string) {
	if v, ok := x.M.LoadAndDelete(key); ok {
		v.(*cron.Cron).Stop()
	}
}

func (x *App) GetState(key string) []cron.Entry {
	if c, ok := x.Get(key); ok {
		return c.Entries()
	}
	return []cron.Entry{}
}

func (x *App) State() (err error) {
	subj := fmt.Sprintf(`%s.schedules-state.%s`, x.V.Namespace, x.V.Id)
	queue := fmt.Sprintf(`%s:schedules-state:%s`, x.V.Namespace, x.V.Id)
	if _, err = x.Nats.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		key := string(msg.Data)
		var states []typ.ScheduleState
		for _, entry := range x.GetState(key) {
			states = append(states, typ.ScheduleState{
				Next: entry.Next,
				Prev: entry.Prev,
			})
		}
		var b []byte
		if b, err = msgpack.Marshal(states); err != nil {
			return
		}
		msg.Respond(b)
	}); err != nil {
		return
	}
	x.Log.Debug("State:ok",
		zap.String("subj", subj),
	)
	return
}
