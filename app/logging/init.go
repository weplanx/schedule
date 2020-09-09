package logging

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"schedule-microservice/app/types"
	pb "schedule-microservice/transfer"
)

type Logging struct {
	types.LoggingOption
	transfer pb.RouterClient
}

func NewLogging(option types.LoggingOption) *Logging {
	c := new(Logging)
	c.LoggingOption = option
	if c.Transfer.Listen != "" {
		conn, err := grpc.Dial(c.Transfer.Listen, grpc.WithInsecure())
		if err != nil {
			logrus.Fatalln(err)
		}
		c.transfer = pb.NewRouterClient(conn)
	}
	return c
}

func (c *Logging) Push(push *types.LoggingPush) (err error) {
	if c.Storage != "" {
		err = c.forStorage(push)
		if err != nil {
			return
		}
	}
	if c.transfer != nil {
		err = c.forTransfer(push)
		if err != nil {
			return
		}
	}
	return
}
