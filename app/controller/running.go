package controller

import (
	"context"
	pb "schedule-microservice/router"
)

func (c *controller) Running(ctx context.Context, param *pb.RunningParameter) (*pb.Response, error) {
	err := c.manager.Running(param.Identity, param.Running)
	if err != nil {
		return c.response(err)
	}
	return c.response(nil)
}
