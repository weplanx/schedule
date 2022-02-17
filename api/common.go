package api

import (
	"github.com/google/wire"
	zlogging "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/weplanx/schedule/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var Provides = wire.NewSet(New)

func New(i *common.Inject) (s *grpc.Server, err error) {
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
	RegisterAPIServer(s, &API{Inject: i})
	return
}
