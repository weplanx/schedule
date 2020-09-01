package controller

import (
	"schedule-microservice/app/manage"
	pb "schedule-microservice/router"
)

type controller struct {
	pb.UnimplementedRouterServer
	manager *manage.JobsManager
}

func New(manager *manage.JobsManager) *controller {
	c := new(controller)
	c.manager = manager
	return c
}

func (c *controller) response(err error) (*pb.Response, error) {
	if err != nil {
		return &pb.Response{
			Error: 1,
			Msg:   err.Error(),
		}, nil
	} else {
		return &pb.Response{
			Error: 0,
			Msg:   "ok",
		}, nil
	}
}
