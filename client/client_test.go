package client_test

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/schedule/client"
	"github.com/weplanx/schedule/common"
	"os"
	"testing"
)

var x *client.Client

func TestMain(m *testing.M) {
	ctx := context.TODO()
	nc, js, err := UseNats(ctx)
	if err != nil {
		panic(err)
	}
	if x, err = client.New(
		client.SetNamespace("example"),
		client.SetNats(nc),
		client.SetJetStream(js),
		client.SetNode(os.Getenv("NODE")),
	); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func UseNats(ctx context.Context) (nc *nats.Conn, js nats.JetStreamContext, err error) {
	var auth nats.Option
	var kp nkeys.KeyPair
	if kp, err = nkeys.FromSeed([]byte(os.Getenv("NATS_NKEY"))); err != nil {
		return
	}
	defer kp.Wipe()
	var pub string
	if pub, err = kp.PublicKey(); err != nil {
		return
	}
	if !nkeys.IsValidPublicUserKey(pub) {
		panic("nkey failed")
	}
	auth = nats.Nkey(pub, func(nonce []byte) ([]byte, error) {
		sig, _ := kp.Sign(nonce)
		return sig, nil
	})
	if nc, err = nats.Connect(
		os.Getenv("NATS_HOSTS"),
		auth,
	); err != nil {
		return
	}
	if js, err = nc.JetStream(
		nats.PublishAsyncMaxPending(256),
		nats.Context(ctx),
	); err != nil {
		return
	}
	return
}

func TestSchedule_Set(t *testing.T) {
	err := x.Set("api", common.ScheduleOption{
		Status: false,
		Jobs: []common.ScheduleJob{
			{
				Mode: "HTTP",
				Spec: "*/2 * * * * *",
				Option: common.HttpOption{
					Url: "https://api.kainonly.com/whoami",
				},
			},
		},
	})
	assert.NoError(t, err)
}

func TestSchedule_Ping(t *testing.T) {
	r, err := x.Ping()
	assert.NoError(t, err)
	t.Log(r)
}

func TestSchedule_Lists(t *testing.T) {
	keys, err := x.Lists()
	assert.NoError(t, err)
	t.Log(keys)
}

func TestSchedule_Get(t *testing.T) {
	jobs, err := x.Get("api")
	assert.NoError(t, err)
	t.Log(jobs)
}

func TestSchedule_StatusStart(t *testing.T) {
	err := x.Status("api", true)
	assert.NoError(t, err)
}

func TestSchedule_StatusStop(t *testing.T) {
	err := x.Status("api", false)
	assert.NoError(t, err)
}

func TestSchedule_Remove(t *testing.T) {
	err := x.Remove("api")
	assert.NoError(t, err)
}
