package common

import (
	"go.uber.org/fx"
	"schedule-microservice/application/service/job"
	"schedule-microservice/application/service/schema"
	"schedule-microservice/config"
)

type Dependency struct {
	fx.In

	Config *config.Config
	Schema *schema.Schema
	Job    *job.Job
}
