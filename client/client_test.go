package client_test

import (
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
	node := os.Getenv("NODE")
	nc, err := UseNats()
	if err != nil {
		panic(err)
	}
	if x, err = client.New(node, nc); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func UseNats() (nc *nats.Conn, err error) {
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
	return
}

func TestSchedule_Set(t *testing.T) {
	err := x.Set("dogs", common.ScheduleOption{
		Status: false,
		Jobs: []common.ScheduleJob{
			{
				Mode: "HTTP",
				Spec: "*/5 * * * * *",
				Option: common.HttpOption{
					Method: "GET",
					Url:    "https://dog.ceo/api/breeds/image/random",
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
	jobs, err := x.Get("dogs")
	assert.NoError(t, err)
	t.Log(jobs)
}

func TestSchedule_StatusStart(t *testing.T) {
	err := x.Status("dogs", true)
	assert.NoError(t, err)
}

func TestSchedule_StatusStop(t *testing.T) {
	err := x.Status("dogs", false)
	assert.NoError(t, err)
}

func TestSchedule_Remove(t *testing.T) {
	err := x.Remove("dogs")
	assert.NoError(t, err)
}
