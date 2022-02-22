package api

import (
	"context"
	"github.com/google/wire"
	zlogging "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/weplanx/schedule/common"
	"github.com/weplanx/schedule/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var Provides = wire.NewSet(New)

func New(i *common.Inject) (s *grpc.Server, err error) {
	ctx := context.Background()

	if err = i.Transfer.CreateLogger(ctx,
		"schedule",
		"schedule",
		"定时器",
	); err != nil {
		return
	}

	coll := i.Values.Database.Collection

	// 初始化存储索引
	if _, err = i.Db.Collection(coll).Indexes().
		CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.M{"key": 1},
				Options: options.Index().
					SetName("uk_key").
					SetUnique(true),
			},
		}); err != nil {
		return
	}

	i.Log.Info("已初始化存储索引",
		zap.Any("state", []interface{}{"uk_key"}),
	)

	// 中间件
	zlogging.ReplaceGrpcLoggerV2(i.Log)
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			zlogging.UnaryServerInterceptor(i.Log),
			recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			zlogging.StreamServerInterceptor(i.Log),
			recovery.StreamServerInterceptor(),
		),
	}

	// 设置 TLS
	if i.Values.TLS.Cert != "" && i.Values.TLS.Key != "" {
		tls := i.Values.TLS
		var creds credentials.TransportCredentials
		if creds, err = credentials.NewServerTLSFromFile(
			tls.Cert,
			tls.Key,
		); err != nil {
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s = grpc.NewServer(opts...)
	api := &API{Inject: i}

	// 启动任务
	var cursor *mongo.Cursor
	if cursor, err = i.Db.Collection(coll).
		Find(ctx, bson.M{
			"node": i.Values.Node,
		}); err != nil {
		return
	}
	schedules := make([]model.Schedule, 0)
	if err = cursor.All(ctx, &schedules); err != nil {
		return
	}
	for _, v := range schedules {
		if err = api.SetSchedule(ctx, v.Key, v.Jobs); err != nil {
			return
		}
	}

	i.Log.Info("任务已启动",
		zap.Any("stats", len(schedules)),
	)
	i.Log.Debug("任务详情",
		zap.Any("schedules", schedules),
	)

	RegisterAPIServer(s, api)
	return
}
