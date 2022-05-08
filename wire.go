//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/weplanx/schedule/app"
	"github.com/weplanx/schedule/bootstrap"
	"github.com/weplanx/schedule/common"
)

func App(value *common.Values) (*app.App, error) {
	wire.Build(
		wire.Struct(new(common.Inject), "*"),
		bootstrap.Provides,
		app.Provides,
	)
	return &app.App{}, nil
}
