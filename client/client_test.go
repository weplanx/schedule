package client

import (
	"context"
	"fmt"
	"github.com/weplanx/schedule/bootstrap"
	"github.com/weplanx/schedule/common"
	"github.com/weplanx/schedule/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"testing"
)

var v *common.Values
var x *Schedule

func TestMain(m *testing.M) {
	os.Chdir("../")
	var err error
	if v, err = bootstrap.SetValues(); err != nil {
		panic(err)
	}
	var host string
	var opts []grpc.DialOption
	if v.TLS.Cert != "" {
		creds, err := credentials.NewClientTLSFromFile(v.TLS.Cert, "")
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
		host = "x.kainonly.com"
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		host = "127.0.0.1"
	}

	if x, err = New(fmt.Sprintf(`%s%s`, host, v.Address), opts...); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestSchedule_Put(t *testing.T) {
	httpJob, err := HttpJob("@every 1s", model.HttpJob{
		Url: "https://api.kainonly.com",
	})
	if err != nil {
		t.Error(err)
	}
	if err = x.Put(context.TODO(), "alpha", httpJob); err != nil {
		t.Error(err)
	}
}

func TestSchedule_Get(t *testing.T) {
	data, err := x.Get(context.TODO(), []string{"alpha"})
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestSchedule_Delete(t *testing.T) {
	if err := x.Delete(context.TODO(), "alpha"); err != nil {
		t.Error(err)
	}
}
