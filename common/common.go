package common

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Inject struct {
	Values   *Values
	Log      *zap.Logger
	Mongo    *mongo.Client
	Db       *mongo.Database
	Schedule *Schedule
}

type Values struct {
	Address  string   `yaml:"address"`
	TLS      TLS      `yaml:"tls"`
	Node     string   `yaml:"node"`
	Database Database `yaml:"database"`
}

type TLS struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type Database struct {
	Uri        string `yaml:"uri"`
	Name       string `yaml:"name"`
	Collection string `yaml:"collection"`
}

func PointInt64(v int64) *int64 {
	return &v
}
