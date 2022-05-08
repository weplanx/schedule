package app

import (
	"github.com/google/wire"
	"github.com/weplanx/schedule/common"
	"github.com/weplanx/schedule/utiliy"
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
