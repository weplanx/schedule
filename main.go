package main

import (
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"schedule-microservice/common"
	"schedule-microservice/controller"
	pb "schedule-microservice/router"
	"schedule-microservice/task"
)

func main() {
	if _, err := os.Stat("./config/autoload"); os.IsNotExist(err) {
		os.Mkdir("./config/autoload", os.ModeDir)
	}
	in, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	cfg := common.AppOption{}
	err = yaml.Unmarshal(in, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	if cfg.Debug {
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}
	tk := task.Create()
	listen, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		log.Fatalln(err)
	}
	server := grpc.NewServer()
	pb.RegisterRouterServer(
		server,
		controller.New(tk),
	)
	server.Serve(listen)
}
