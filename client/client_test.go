package client

import (
	"github.com/weplanx/schedule/bootstrap"
	"github.com/weplanx/schedule/common"
	"os"
	"testing"
)

var x *Schedule

func TestMain(m *testing.M) {
	os.Chdir("../")
	values, err := common.SetValues()
	if err != nil {
		panic(err)
	}
	nc, err := bootstrap.UseNats(values)
	if err != nil {
		panic(err)
	}
	js, err := bootstrap.UseJetStream(nc)
	if err != nil {
		panic(err)
	}
	if x, err = New(values.Namespace, nc, js); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestSchedule_List(t *testing.T) {
	result, err := x.List()
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

var key = "69dd963c-2ec3-1f11-f1a5-ddcc1daf1f8f"

func TestSchedule_Set(t *testing.T) {
	job1 := common.HttpJob("@every 5s", common.HttpOption{
		Url: "https://api.kainonly.com",
	})
	job2 := common.HttpJob("@every 10s", common.HttpOption{
		Url: "https://api.kainonly.com",
	})
	if err := x.Set(key, job1, job2); err != nil {
		t.Error(err)
	}
}

func TestSchedule_Update(t *testing.T) {
	job := common.HttpJob("@every 15s", common.HttpOption{
		Url: "https://api.kainonly.com",
		Headers: map[string]string{
			"x-token": "6xvLvuQUhc2$j8H#",
		},
		Body: map[string]interface{}{
			"name": key,
		},
	})
	if err := x.Set(key, job); err != nil {
		t.Error(err)
	}
}

func TestSchedule_Get(t *testing.T) {
	data, err := x.Get(key)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestSchedule_Status(t *testing.T) {
	result, err := x.Status("62ad37d1cb8a8fb377bccae1", false)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(result))
}

func TestSchedule_Remove(t *testing.T) {
	if err := x.Remove(key); err != nil {
		t.Error(err)
	}
}
