//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/workflow/worker/app"
	"github.com/weplanx/workflow/worker/common"
)

func NewApp() (*app.App, error) {
	wire.Build(
		wire.Struct(new(common.Inject), "*"),
		LoadStaticValues,
		UseZap,
		UseNats,
		UseJetStream,
		UseTransfer,
		app.Initialize,
	)
	return &app.App{}, nil
}
