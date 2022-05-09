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
	if x, err = New("alpha", nc, js); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestSchedule_Set(t *testing.T) {
	job := common.HttpJob("@every 1s", common.HttpOption{
		Url: "http://mac:8080/ping",
	})
	if err := x.Set("ping", job); err != nil {
		t.Error(err)
	}
}

func TestSchedule_Update(t *testing.T) {
	job := common.HttpJob("@every 5s", common.HttpOption{
		Url: "http://mac:8080/ping",
		Headers: map[string]string{
			"x-token": "zxc",
		},
		Body: map[string]interface{}{
			"name": "kain",
		},
	})
	if err := x.Set("ping", job); err != nil {
		t.Error(err)
	}
}

func TestSchedule_Get(t *testing.T) {
	data, err := x.Get("ping")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestSchedule_Remove(t *testing.T) {
	if err := x.Remove("ping"); err != nil {
		t.Error(err)
	}
}
