package app

import (
	"google.golang.org/grpc"
	"net"
	"net/http"
	_ "net/http/pprof"
	"schedule-microservice/app/controller"
	"schedule-microservice/app/logging"
	"schedule-microservice/app/manage"
	"schedule-microservice/app/schema"
	"schedule-microservice/app/types"
	pb "schedule-microservice/router"
)

func Application(option *types.Config) (err error) {
	// Turn on debugging
	if option.Debug {
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}
	// Start microservice
	listen, err := net.Listen("tcp", option.Listen)
	if err != nil {
		return
	}
	server := grpc.NewServer()
	dataset := schema.New()
	logger := logging.NewLogging(option.Logging)
	manager, err := manage.NewJobsManager(dataset, logger)
	if err != nil {
		return
	}
	pb.RegisterRouterServer(
		server,
		controller.New(manager),
	)
	server.Serve(listen)
	return
}
