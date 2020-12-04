package transfer

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"reflect"
	"schedule-microservice/config/options"
	pb "schedule-microservice/transfer"
)

type Transfer struct {
	on     bool
	client pb.APIClient
	pipe   string
}

func New(option options.TransferOption) (c *Transfer, err error) {
	c = new(Transfer)
	if reflect.DeepEqual(option, options.TransferOption{}) {
		c.on = false
		return
	}
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(option.Listen, grpc.WithInsecure()); err != nil {
		return
	}
	c.on = true
	c.client = pb.NewAPIClient(conn)
	c.pipe = option.Pipe
	return
}

func (c *Transfer) Push(value interface{}) (err error) {
	if !c.on {
		return
	}
	var data []byte
	if data, err = jsoniter.Marshal(value); err != nil {
		return
	}
	if _, err = c.client.Push(context.Background(), &pb.Body{
		Id:      c.pipe,
		Content: data,
	}); err != nil {
		return
	}
	return
}
