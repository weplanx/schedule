package app

import (
	"fmt"
	"github.com/google/wire"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/common"
	"github.com/weplanx/schedule/utiliy"
	"time"
)

var Provides = wire.NewSet(New)

type App struct {
	*common.Inject
	Schedules *utiliy.Schedule
}

func New(i *common.Inject) (x *App, err error) {
	x = &App{
		Inject:    i,
		Schedules: utiliy.NewSchedule(),
	}
	return
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
