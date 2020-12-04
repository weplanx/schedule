package controller

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	pb "schedule-microservice/api"
)

func (c *controller) Running(_ context.Context, option *pb.Status) (*empty.Empty, error) {
	if err := c.Job.Running(option.Id, option.Running); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
