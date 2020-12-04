package controller

import (
	"context"
	pb "schedule-microservice/api"
)

func (c *controller) Lists(_ context.Context, option *pb.IDs) (_ *pb.Options, err error) {
	lists := make([]*pb.Option, len(option.Ids))
	for key, val := range option.Ids {
		if lists[key], err = c.find(val); err != nil {
			return
		}
	}
	return &pb.Options{Options: lists}, nil
}
