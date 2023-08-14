package workflow_test

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/weplanx/workflow"
	"os"
	"testing"
)

var x *workflow.Workflow
var schedule *workflow.Schedule

func TestMain(m *testing.M) {
	ctx := context.TODO()
	nc, js, err := UseNats(ctx)
	if err != nil {
		panic(err)
	}
	x = workflow.New(
		workflow.SetNamespace("example"),
		workflow.SetNats(nc),
		workflow.SetJetStream(js),
	)
	if schedule, err = x.NewSchedule(os.Getenv("ID")); err != nil {
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
