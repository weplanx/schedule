package bootstrap

import (
	"context"
	"errors"
	"github.com/google/wire"
	"github.com/weplanx/schedule/common"
	"github.com/weplanx/transfer/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var Provides = wire.NewSet(
	UseZap,
	UseMongoDB,
	UseDatabase,
	UseSchedule,
	UseTransfer,
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

func UseSchedule() *common.Schedule {
	return common.NewSchedule()
}

func UseTransfer(values *common.Values) (*client.Transfer, error) {
	option := values.Transfer
	var opts []grpc.DialOption
	if option.TLS.Cert != "" {
		creds, err := credentials.NewClientTLSFromFile(option.TLS.Cert, "")
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return client.New(option.Address, opts...)
}
