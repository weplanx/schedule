//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/schedule/app"
	"github.com/weplanx/schedule/common"
)

func NewApp() (*app.App, error) {
	wire.Build(
		wire.Struct(new(common.Inject), "*"),
		Provides,
		app.Provides,
	)
	return &app.App{}, nil
}
