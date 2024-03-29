package app

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/common"
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
		var option common.ScheduleOption
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
				var option common.ScheduleOption
				if err = msgpack.Unmarshal(update.Value(), &option); err != nil {
					return
				}
				if err = x.Set(key, option); err != nil {
					x.Log.Error("schedule fail",
						zap.String("key", key),
						zap.Any("option", option),
						zap.Error(err),
					)
					return
				}
				x.Log.Debug("schedule ok",
					zap.String("key", key),
					zap.Any("option", option),
				)
				break
			case "KeyValueDeleteOp":
				x.Remove(key)
				x.Log.Debug("schedule removed",
					zap.String("key", key),
				)
				break
			}
		}
	}()

	if err = x.State(); err != nil {
		return
	}
	if err = x.Ping(); err != nil {
		return
	}
	x.Log.Info("service started!")
	return
}

func (x *App) Get(key string) (*cron.Cron, bool) {
	v, ok := x.M.Load(key)
	return v.(*cron.Cron), ok
}

func (x *App) Set(key string, option common.ScheduleOption) (err error) {
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

func (x *App) SetJob(key string, c *cron.Cron, index int, job common.ScheduleJob) (err error) {
	if _, err = c.AddFunc(job.Spec, func() {
		subj := fmt.Sprintf(`jobs.%s`, x.V.Node)
		var b []byte
		if b, err = msgpack.Marshal(common.Job{
			Key:    key,
			Index:  index,
			Mode:   job.Mode,
			Option: job.Option,
		}); err != nil {
			return
		}
		if err = x.Nats.Publish(subj, b); err != nil {
			x.Log.Error("publish fail",
				zap.String("key", key),
				zap.Int("index", index),
				zap.Error(err),
			)
			return
		}
		x.Log.Debug("publish ok",
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
	subj := fmt.Sprintf(`schedules.%s`, x.V.Node)
	queue := fmt.Sprintf(`SCHEDULE_%s`, x.V.Node)
	if _, err = x.Nats.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		key := string(msg.Data)
		var states []common.ScheduleState
		for _, entry := range x.GetState(key) {
			states = append(states, common.ScheduleState{
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
	x.Log.Debug("state ok",
		zap.String("subj", subj),
	)
	return
}

func (x *App) Ping() (err error) {
	if _, err = x.Nats.Subscribe("schedules", func(msg *nats.Msg) {
		if string(msg.Data) == x.V.Node {
			msg.Respond([]byte("ok"))
		}
	}); err != nil {
		return
	}
	return
}
