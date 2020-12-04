package controller

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	pb "schedule-microservice/api"
)

func (c *controller) Running(_ context.Context, option *pb.Status) (*empty.Empty, error) {
	//err := c.manager.Running(param.Identity, param.Running)
	//if err != nil {
	//	return c.response(err)
	//}
	return nil, nil
}
