package main

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"schedule-microservice/common"
	"schedule-microservice/controller"
	"schedule-microservice/job"
	pb "schedule-microservice/router"
)

func main() {
	if _, err := os.Stat("./config/autoload"); os.IsNotExist(err) {
		os.Mkdir("./config/autoload", os.ModeDir)
	}
	in, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		logrus.Fatalln(err)
	}
	cfg := common.AppOption{}
	err = yaml.Unmarshal(in, &cfg)
	if err != nil {
		logrus.Fatalln(err)
	}
	if cfg.Debug {
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}
	err = common.SetLogger(&cfg.Log)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer common.LoggerClose()
	comJob := job.Create()
	listen, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		logrus.Fatalln(err)
	}
	server := grpc.NewServer()
	pb.RegisterRouterServer(
		server,
		controller.New(comJob),
	)
	server.Serve(listen)
}
