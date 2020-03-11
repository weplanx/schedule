package controller

import (
	"context"
	pb "schedule-microservice/router"
)

type controller struct {
	pb.UnimplementedRouterServer
}

func New() *controller {
	c := new(controller)
	return c
}

func (c *controller) Get(ctx context.Context, req *pb.GetParameter) (*pb.GetResponse, error) {
	return &pb.GetResponse{
		Error: 0,
		Msg:   "",
		Data:  nil,
	}, nil
}

func (c *controller) Lists(ctx context.Context, req *pb.ListsParameter) (*pb.ListsResponse, error) {
	return &pb.ListsResponse{
		Error: 0,
		Msg:   "",
		Data:  nil,
	}, nil
}

func (c *controller) All(ctx context.Context, req *pb.NoParameter) (*pb.AllResponse, error) {
	return &pb.AllResponse{
		Error: 0,
		Msg:   "",
		Data:  nil,
	}, nil
}

func (c *controller) Put(ctx context.Context, req *pb.PutParameter) (*pb.Response, error) {
	return &pb.Response{
		Error: 0,
		Msg:   "",
	}, nil
}

func (c *controller) Delete(ctx context.Context, req *pb.DeleteParameter) (*pb.Response, error) {
	return &pb.Response{
		Error: 0,
		Msg:   "",
	}, nil
}

func (c *controller) Running(ctx context.Context, req *pb.RunningParameter) (*pb.Response, error) {
	return &pb.Response{
		Error: 0,
		Msg:   "",
	}, nil
}
