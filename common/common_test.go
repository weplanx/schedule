package common

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Chdir("..")
	os.Exit(m.Run())
}

func TestHttpClient(t *testing.T) {
	agent := gorequest.New().Post("http://localhost:3000")
	agent.Set("X-TOKEN", "vvv")
	agent.Send(`{"order":"x-x1"}`)
	_, body, errs := agent.EndBytes()
	if errs != nil {
		t.Fatal(errs)
	} else {
		println(string(body))
	}
}

func TestConfig(t *testing.T) {
	if _, err := os.Stat("./config/autoload"); os.IsNotExist(err) {
		os.Mkdir("./config/autoload", os.ModeDir)
	}
}

func TestSaveConfig(t *testing.T) {
	var body1 interface{}
	err := json.Unmarshal([]byte(`{"name":"task1"}`), &body1)
	if err != nil {
		t.Fatal(err)
	}
	var body2 interface{}
	err = json.Unmarshal([]byte(`{"name":"task2"}`), &body2)
	if err != nil {
		t.Fatal(err)
	}
	data := &TaskOption{
		Identity: "test",
		TimeZone: "Asia/Shanghai",
		Start:    true,
		Entries: map[string]*EntryOption{
			"task1": &EntryOption{
				CronTime: "*/10 * * * * *",
				Url:      "http://localhost:3000",
				Headers: map[string]string{
					"x-token": "abc",
				},
				Body: body1,
			},
			"task2": &EntryOption{
				CronTime: "*/20 * * * * *",
				Url:      "http://localhost:3000",
				Headers: map[string]string{
					"x-token": "abc",
				},
				Body: body2,
			},
		},
	}
	err = SaveConfig(data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListConfig(t *testing.T) {
	lists, err := ListConfig()
	if err != nil {
		t.Fatal(err)
	}
	logrus.Info(lists)
}

func TestRemoveConfig(t *testing.T) {
	err := RemoveConfig("test")
	if err != nil {
		t.Fatal(err)
	}
}
