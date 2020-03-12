package controller

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"schedule-microservice/common"
	pb "schedule-microservice/router"
	"testing"
)

var conn *grpc.ClientConn

func TestMain(m *testing.M) {
	os.Chdir("..")
	in, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	cfg := common.AppOption{}
	err = yaml.Unmarshal(in, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err = grpc.Dial(cfg.Listen, grpc.WithInsecure())
	os.Exit(m.Run())
}

func TestPut(t *testing.T) {
	defer conn.Close()
	client := pb.NewRouterClient(conn)
	response, err := client.Put(
		context.Background(),
		&pb.PutParameter{
			Identity: "test",
			TimeZone: "Asia/Shanghai",
			Start:    true,
			Entries: map[string]*pb.EntryOption{
				"task1": &pb.EntryOption{
					CronTime: "*/5 * * * * *",
					Url:      "http://localhost:3000",
					Headers:  []byte(`{"x-token":"abc"}`),
					Body:     []byte(`{"name":"task1"}`),
				},
				"task2": &pb.EntryOption{
					CronTime: "*/10 * * * * *",
					Url:      "http://localhost:3000",
					Headers:  []byte(`{"x-token":"abc"}`),
					Body:     []byte(`{"name":"task2"}`),
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error == 1 {
		t.Fatal(response.Msg)
	}
}

func TestPutOther(t *testing.T) {
	defer conn.Close()
	client := pb.NewRouterClient(conn)
	response, err := client.Put(
		context.Background(),
		&pb.PutParameter{
			Identity: "next",
			TimeZone: "Asia/Shanghai",
			Start:    true,
			Entries: map[string]*pb.EntryOption{
				"other1": &pb.EntryOption{
					CronTime: "*/8 * * * * *",
					Url:      "http://localhost:3000",
					Headers:  []byte(`{"x-token":"123"}`),
					Body:     []byte(`{"name":"other1"}`),
				},
				"other2": &pb.EntryOption{
					CronTime: "*/16 * * * * *",
					Url:      "http://localhost:3000",
					Headers:  []byte(`{"x-token":"123"}`),
					Body:     []byte(`{"name":"other2"}`),
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error == 1 {
		t.Fatal(response.Msg)
	}
}

func TestGet(t *testing.T) {
	defer conn.Close()
	client := pb.NewRouterClient(conn)
	response, err := client.Get(
		context.Background(),
		&pb.GetParameter{
			Identity: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error == 1 {
		t.Fatal(response.Msg)
	}
	logrus.Info(response.Data)
}

func TestLists(t *testing.T) {
	defer conn.Close()
	client := pb.NewRouterClient(conn)
	response, err := client.Lists(
		context.Background(),
		&pb.ListsParameter{
			Identity: []string{"test", "other"},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error == 1 {
		t.Fatal(response.Msg)
	}
	logrus.Info(response.Data)
}

func TestAll(t *testing.T) {
	defer conn.Close()
	client := pb.NewRouterClient(conn)
	response, err := client.All(
		context.Background(),
		&pb.NoParameter{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error == 1 {
		t.Fatal(response.Msg)
	}
	logrus.Info(response.Data)
}

func TestDelete(t *testing.T) {
	defer conn.Close()
	client := pb.NewRouterClient(conn)
	response, err := client.Delete(
		context.Background(),
		&pb.DeleteParameter{
			Identity: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error == 1 {
		t.Fatal(response.Msg)
	}
}

func TestRunning(t *testing.T) {
	defer conn.Close()
	client := pb.NewRouterClient(conn)
	response, err := client.Running(
		context.Background(),
		&pb.RunningParameter{
			Identity: "test",
			Running:  false,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error == 1 {
		t.Fatal(response.Msg)
	}
}
