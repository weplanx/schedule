//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/workflow/excel/api"
	"github.com/weplanx/workflow/excel/common"
)

func NewAPI() (*api.API, error) {
	wire.Build(
		wire.Struct(new(api.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		LoadValues,
		UseCos,
	)
	return &api.API{}, nil
}
