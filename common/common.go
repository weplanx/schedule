package common

import (
	"github.com/weplanx/schedule/app"
	"github.com/weplanx/transfer/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Inject struct {
	Values   *Values
	Log      *zap.Logger
	Mongo    *mongo.Client
	Db       *mongo.Database
	Schedule *app.Schedule
	Transfer *client.Transfer
}

type Values struct {
	Address  string   `yaml:"address"`
	TLS      TLS      `yaml:"tls"`
	Node     string   `yaml:"node"`
	Database Database `yaml:"database"`
	Transfer Transfer `yaml:"transfer"`
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

type Transfer struct {
	Address string `yaml:"address"`
	TLS     TLS    `yaml:"tls"`
}
