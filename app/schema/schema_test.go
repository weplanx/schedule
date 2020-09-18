package schema

import (
	jsoniter "github.com/json-iterator/go"
	"os"
	"schedule-microservice/app/types"
	"testing"
)

var schema *Schema

func TestMain(m *testing.M) {
	os.Chdir("../..")
	schema = New()
	os.Exit(m.Run())
}

func TestSchema_Update(t *testing.T) {
	var err error
	var body1 interface{}
	err = jsoniter.Unmarshal([]byte(`{"name":"task1"}`), &body1)
	if err != nil {
		t.Fatal(err)
	}
	var body2 interface{}
	err = jsoniter.Unmarshal([]byte(`{"name":"task2"}`), &body2)
	if err != nil {
		t.Fatal(err)
	}
	err = schema.Update(types.JobOption{
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
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSchema_Lists(t *testing.T) {
	_, err := schema.Lists()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSchema_Delete(t *testing.T) {
	err := schema.Delete("test")
	if err != nil {
		t.Fatal(err)
	}
}
