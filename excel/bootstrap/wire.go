//go:build wireinject
// +build wireinject

package bootstrap

import (
	"excel/api"
	"excel/common"
	"github.com/google/wire"
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
