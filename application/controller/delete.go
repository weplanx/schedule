package controller

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	pb "schedule-microservice/api"
)

func (c *controller) Delete(_ context.Context, option *pb.ID) (*empty.Empty, error) {
	return nil, nil
}
