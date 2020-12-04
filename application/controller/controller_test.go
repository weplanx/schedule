package controller

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"os"
	pb "schedule-microservice/api"
	"schedule-microservice/bootstrap"
	"schedule-microservice/config"
	"testing"
)

var client pb.APIClient

func TestMain(m *testing.M) {
	os.Chdir("../../")
	var err error
	var cfg *config.Config
	if cfg, err = bootstrap.LoadConfiguration(); err != nil {
		log.Fatalln(err)
	}
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(cfg.Listen, grpc.WithInsecure()); err != nil {
		log.Fatalln(err)
	}
	client = pb.NewAPIClient(conn)
	os.Exit(m.Run())
}

func TestController_Put(t *testing.T) {
	_, err := client.Put(context.Background(), &pb.Option{
		Id:       "debug-A",
		TimeZone: "Asia/Shanghai",
		Start:    true,
		Entries: map[string]*pb.Entry{
			"entry-1": {
				CronTime: "*/10 * * * * *",
				Url:      "http://mac:3000/entry-1",
				Headers:  []byte(`{"x-token":"l51aM51gp43606o2"}`),
				Body:     []byte(`{"msg":"hello entry-A1"}`),
			},
			"entry-2": {
				CronTime: "*/20 * * * * *",
				Url:      "http://mac:3000/entry-2",
				Headers:  []byte(`{"x-token":"GGlxNXfMyJb5IKuL"}`),
				Body:     []byte(`{"msg":"hello entry-A2"}`),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Put(context.Background(), &pb.Option{
		Id:       "debug-B",
		TimeZone: "Asia/Shanghai",
		Start:    true,
		Entries: map[string]*pb.Entry{
			"entry-1": {
				CronTime: "*/5 * * * * *",
				Url:      "http://mac:3000/task3",
				Headers:  []byte(`{"x-token":"ymNS2ZZzKKbqWpVm"}`),
				Body:     []byte(`{"msg":"hello entry-B1"}`),
			},
			"entry-2": {
				CronTime: "*/30 * * * * *",
				Url:      "http://mac:3000/task4",
				Headers:  []byte(`{"x-token":"AFAghq7Nc8S5gDr4"}`),
				Body:     []byte(`{"msg":"hello entry-B2"}`),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestController_Get(t *testing.T) {
	response, err := client.Get(context.Background(), &pb.ID{
		Id: "debug-A",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}

func TestController_Lists(t *testing.T) {
	response, err := client.Lists(context.Background(), &pb.IDs{
		Ids: []string{"debug-A", "debug-B"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}

func TestController_All(t *testing.T) {
	response, err := client.All(context.Background(), &empty.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}

func TestController_Running(t *testing.T) {
	_, err := client.Running(context.Background(), &pb.Status{
		Id:      "debug-A",
		Running: false,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestController_Delete(t *testing.T) {
	_, err := client.Delete(context.Background(), &pb.ID{
		Id: "debug-A",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Delete(context.Background(), &pb.ID{
		Id: "debug-B",
	})
	if err != nil {
		t.Fatal(err)
	}
}
