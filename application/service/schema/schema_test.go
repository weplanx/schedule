package schema

import (
	jsoniter "github.com/json-iterator/go"
	"os"
	"schedule-microservice/config/options"
	"testing"
)

var schema *Schema

func TestMain(m *testing.M) {
	os.Chdir("../../..")
	schema = New("./config/autoload/")
	os.Exit(m.Run())
}

func TestSchema_Update(t *testing.T) {
	var body interface{}
	if err := jsoniter.Unmarshal([]byte(`{"name":"hello"}`), &body); err != nil {
		t.Fatal(err)
	}
	err := schema.Update(options.JobOption{
		Identity: "debug",
		TimeZone: "Asia/Shanghai",
		Start:    true,
		Entries: map[string]*options.EntryOption{
			"debug": {
				CronTime: "*/10 * * * * *",
				Url:      "http://localhost:3000/hello",
				Headers: map[string]string{
					"x-token": "l51aM51gp43606o2",
				},
				Body: body,
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
	err := schema.Delete("debug")
	if err != nil {
		t.Fatal(err)
	}
}
