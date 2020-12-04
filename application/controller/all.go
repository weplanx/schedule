package controller

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	pb "schedule-microservice/api"
)

func (c *controller) All(_ context.Context, _ *empty.Empty) (*pb.IDs, error) {
	var ids []string
	for id, _ := range c.Job.Options.Lists() {
		ids = append(ids, id)
	}
	return &pb.IDs{Ids: ids}, nil
}
