package common

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Inject struct {
	Values *Values
	Log    *zap.Logger
	Mongo  *mongo.Client
	Db     *mongo.Database
}

type Values struct {
	Address   string   `yaml:"address"`
	TLS       TLS      `yaml:"tls"`
	Namespace string   `yaml:"namespace"`
	Database  Database `yaml:"database"`
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
