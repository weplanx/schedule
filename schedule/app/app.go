package app

import (
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
	now := time.Now()
	var watch nats.KeyWatcher
	if watch, err = x.KeyValue.WatchAll(); err != nil {
		return
	}
	for entry := range watch.Updates() {
		if entry == nil || entry.Created().Unix() < now.Unix() {
			continue
		}
		key := entry.Key()
		switch entry.Operation().String() {
		case "KeyValuePutOp":
			var jobs []typ.ScheduleJob
			msgpack.Unmarshal(entry.Value(), &jobs)
			if err = x.Set(key, jobs...); err != nil {
				x.Log.Error("KeyValuePutOp Set:fail",
					zap.String("key", key),
					zap.Any("jobs", jobs),
					zap.Error(err),
				)
				return
			}
			x.Log.Debug("KeyValuePutOp SetCron:ok",
				zap.String("key", key),
				zap.Any("jobs", jobs),
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
	return
}

func (x *App) Get(key string) (*cron.Cron, bool) {
	v, ok := x.M.Load(key)
	return v.(*cron.Cron), ok
}

func (x *App) Set(key string, jobs ...typ.ScheduleJob) (err error) {
	x.M.Delete(key)
	c := cron.New(cron.WithSeconds())
	for i, job := range jobs {
		if err = x.SetJob(key, c, i, job); err != nil {
			return
		}
	}
	x.M.Store(key, c)
	return
}

func (x *App) SetJob(key string, c *cron.Cron, index int, job typ.ScheduleJob) (err error) {
	if _, err = c.AddFunc(job.Spec, func() {
		subj := fmt.Sprintf(`%s.jobs.%s`, x.V.Namespace, x.V.Id)
		b, _ := msgpack.Marshal(typ.Job{
			Key:    key,
			Index:  index,
			Mode:   job.Mode,
			Option: job.Option,
		})
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

func (x *App) Start(key string) {
	if c, ok := x.Get(key); ok {
		c.Start()
	}
}

func (x *App) Stop(key string) {
	if c, ok := x.Get(key); ok {
		c.Stop()
	}
}

func (x *App) GetState(key string) []cron.Entry {
	if c, ok := x.Get(key); ok {
		return c.Entries()
	}
	return []cron.Entry{}
}

//// State job info event
//func (x *App) State() (err error) {
//	name := fmt.Sprintf(`schedules-state:%s`, x.Values.Id)
//	subject := fmt.Sprintf(`schedules-state.%s`, x.Values.Id)
//	if _, err = x.Nats.QueueSubscribe(subject, name, func(msg *nats.Msg) {
//		key := string(msg.Data)
//		var states []typ.ScheduleState
//		for _, entry := range x.GetState(key) {
//			states = append(states, typ.ScheduleState{
//				Next: entry.Next,
//				Prev: entry.Prev,
//			})
//		}
//		b, _ := msgpack.Marshal(states)
//		msg.Respond(b)
//	}); err != nil {
//		return
//	}
//	return
//}
//
//// Status status event
//func (x *App) Status() (err error) {
//	subject := fmt.Sprintf(`schedules-status.%s`, x.Values.Id)
//	if _, err = x.Nats.Subscribe(subject, func(msg *nats.Msg) {
//		var status typ.ScheduleStatus
//		msgpack.Unmarshal(msg.Data, &status)
//		if status.Running {
//			x.Start(status.Key)
//			x.Log.Debug("status started",
//				zap.String("key", status.Key),
//			)
//		} else {
//			x.Stop(status.Key)
//			x.Log.Debug("status stopped",
//				zap.String("key", status.Key),
//			)
//		}
//		msg.Respond([]byte("ok"))
//	}); err != nil {
//		return
//	}
//	return
//}
