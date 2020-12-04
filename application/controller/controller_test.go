package controller

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"schedule-microservice/app/types"
	pb "schedule-microservice/router"
	"strconv"
	"testing"
)

var client pb.RouterClient

func TestMain(m *testing.M) {
	os.Chdir("../..")
	if _, err := os.Stat("./config/autoload"); os.IsNotExist(err) {
		os.Mkdir("./config/autoload", os.ModeDir)
	}
	if _, err := os.Stat("./config/config.yml"); os.IsNotExist(err) {
		logrus.Fatalln("The service configuration file does not exist")
	}
	cfgByte, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		logrus.Fatalln("Failed to read service configuration file", err)
	}
	config := types.Config{}
	err = yaml.Unmarshal(cfgByte, &config)
	if err != nil {
		logrus.Fatalln("Service configuration file parsing failed", err)
	}
	grpcConn, err := grpc.Dial(config.Listen, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalln(err)
	}
	client = pb.NewRouterClient(grpcConn)
	os.Exit(m.Run())
}

func TestController_Put(t *testing.T) {
	response, err := client.Put(context.Background(), &pb.PutParameter{
		Identity: "test-1",
		TimeZone: "Asia/Shanghai",
		Start:    true,
		Entries: map[string]*pb.EntryOption{
			"task1": {
				CronTime: "*/10 * * * * *",
				Url:      "http://mac:3000/task1",
				Headers:  []byte(`{"x-token":"l51aM51gp43606o2"}`),
				Body:     []byte(`{"name":"task1"}`),
			},
			"task2": {
				CronTime: "*/20 * * * * *",
				Url:      "http://mac:3000/task2",
				Headers:  []byte(`{"x-token":"GGlxNXfMyJb5IKuL"}`),
				Body:     []byte(`{"name":"task2"}`),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.Error != 0 {
		t.Error(response.Msg)
	} else {
		t.Log(response.Msg)
	}
	response, err = client.Put(context.Background(), &pb.PutParameter{
		Identity: "test-2",
		TimeZone: "Asia/Shanghai",
		Start:    true,
		Entries: map[string]*pb.EntryOption{
			"task1": {
				CronTime: "*/10 * * * * *",
				Url:      "http://mac:3000/task3",
				Headers:  []byte(`{"x-token":"ymNS2ZZzKKbqWpVm"}`),
				Body:     []byte(`{"name":"task3"}`),
			},
			"task2": {
				CronTime: "*/20 * * * * *",
				Url:      "http://mac:3000/task4",
				Headers:  []byte(`{"x-token":"AFAghq7Nc8S5gDr4"}`),
				Body:     []byte(`{"name":"task4"}`),
			},
		},
	})
}

func TestController_Get(t *testing.T) {
	response, err := client.Get(context.Background(), &pb.GetParameter{
		Identity: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.Error != 0 {
		t.Error(response.Msg)
	} else {
		t.Log(response.Data)
	}
}

func TestController_Lists(t *testing.T) {
	response, err := client.Lists(context.Background(), &pb.ListsParameter{
		Identity: []string{"test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.Error != 0 {
		t.Error(response.Msg)
	} else {
		t.Log(response.Data)
	}
}

func TestController_All(t *testing.T) {
	response, err := client.All(context.Background(), &pb.NoParameter{})
	if err != nil {
		t.Fatal(err)
	}
	if response.Error != 0 {
		t.Error(response.Msg)
	} else {
		t.Log(response.Data)
	}
}

func TestController_Running(t *testing.T) {
	response, err := client.Running(context.Background(), &pb.RunningParameter{
		Identity: "test",
		Running:  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.Error != 0 {
		t.Error(response.Msg)
	} else {
		t.Log(response.Msg)
	}
}

func TestController_Delete(t *testing.T) {
	response, err := client.Delete(context.Background(), &pb.DeleteParameter{
		Identity: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.Error != 0 {
		t.Error(response.Msg)
	} else {
		t.Log(response.Msg)
	}
}

func BenchmarkController_Put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		response, err := client.Put(context.Background(), &pb.PutParameter{
			Identity: "test-" + strconv.Itoa(i),
			TimeZone: "Asia/Shanghai",
			Start:    true,
			Entries: map[string]*pb.EntryOption{
				"task1": {
					CronTime: "*/10 * * * * *",
					Url:      "http://mac:3000/task1",
					Headers:  []byte(`{"x-token":"l51aM51gp43606o2"}`),
					Body:     []byte(`{"name":"task1"}`),
				},
				"task2": {
					CronTime: "*/20 * * * * *",
					Url:      "http://mac:3000/task2",
					Headers:  []byte(`{"x-token":"GGlxNXfMyJb5IKuL"}`),
					Body:     []byte(`{"name":"task2"}`),
				},
			},
		})
		if err != nil {
			b.Fatal(err)
		}
		if response.Error != 0 {
			b.Error(response.Msg)
		}
	}
}
