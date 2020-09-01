package manage

import (
	"encoding/json"
	"log"
	"os"
	"schedule-microservice/app/types"
	"testing"
)

var manager *JobsManager
var option types.JobOption

func TestMain(m *testing.M) {
	os.Chdir("../..")
	var err error
	manager, err = NewJobsManager()
	if err != nil {
		log.Fatalln(err)
	}
	var body1 interface{}
	err = json.Unmarshal([]byte(`{"name":"task1"}`), &body1)
	if err != nil {
		log.Fatalln(err)
	}
	var body2 interface{}
	err = json.Unmarshal([]byte(`{"name":"task2"}`), &body2)
	if err != nil {
		log.Fatalln(err)
	}
	option = types.JobOption{
		Identity: "test",
		TimeZone: "Asia/Shanghai",
		Start:    true,
		Entries: map[string]*types.EntryOption{
			"task1": {
				CronTime: "*/10 * * * * *",
				Url:      "http://localhost:3000/task1",
				Headers: map[string]string{
					"x-token": "l51aM51gp43606o2",
				},
				Body: body1,
			},
			"task2": {
				CronTime: "*/20 * * * * *",
				Url:      "http://localhost:3000/task2",
				Headers: map[string]string{
					"x-token": "GGlxNXfMyJb5IKuL",
				},
				Body: body2,
			},
		},
	}
	os.Exit(m.Run())
}

func TestJobsManager_Put(t *testing.T) {
	err := manager.Put(option)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJobsManager_Running(t *testing.T) {
	err := manager.Running("test", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJobsManager_Delete(t *testing.T) {
	err := manager.Delete("test")
	if err != nil {
		t.Fatal(err)
	}
}
