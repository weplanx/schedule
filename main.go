package main

import (
	"go.uber.org/fx"
	"schedule-microservice/application"
	"schedule-microservice/bootstrap"
)

func main() {
	fx.New(
		//fx.NopLogger,
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeSchema,
			bootstrap.InitializeFilelog,
			bootstrap.InitializeTransfer,
			bootstrap.InitializeJob,
		),
		fx.Invoke(application.Application),
	).Run()
}
