package bootstrap

import (
	"context"
	"errors"
	"github.com/google/wire"
	"github.com/weplanx/schedule/app"
	"github.com/weplanx/schedule/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var Provides = wire.NewSet(
	UseZap,
	UseMongoDB,
	UseDatabase,
	UseSchedule,
)

// SetValues 初始化配置
func SetValues() (values *common.Values, err error) {
	if _, err = os.Stat("./config/config.yml"); os.IsNotExist(err) {
		err = errors.New("the path [./config.yml] does not have a configuration file")
		return
	}
	var b []byte
	b, err = ioutil.ReadFile("./config/config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &values)
	if err != nil {
		return
	}
	return
}

func UseZap() (logger *zap.Logger, err error) {
	if logger, err = zap.NewProduction(); err != nil {
		return
	}
	return
}

func UseMongoDB(values *common.Values) (*mongo.Client, error) {
	return mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(values.Database.Uri),
	)
}

func UseDatabase(client *mongo.Client, values *common.Values) *mongo.Database {
	return client.Database(values.Database.Name)
}

func UseSchedule() *app.Schedule {
	return app.NewSchedule()
}
