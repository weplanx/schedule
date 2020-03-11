package controller

import (
	"context"
	"encoding/json"
	pb "schedule-microservice/router"
	"schedule-microservice/task"
)

type controller struct {
	tk *task.Task
	pb.UnimplementedRouterServer
}

func New() *controller {
	c := new(controller)
	return c
}

func (c *controller) Get(ctx context.Context, req *pb.GetParameter) (*pb.GetResponse, error) {
	var err error
	opt := c.tk.Get(req.Identity)
	var entries []*pb.EntryOptionWithTime
	for _, val := range opt.Entries {
		var headers []byte
		headers, err = json.Marshal(val.Headers)
		if err != nil {
			return &pb.GetResponse{
				Error: 1,
				Msg:   err.Error(),
			}, nil
		}
		var body []byte
		body, err = json.Marshal(val.Body)
		if err != nil {
			return &pb.GetResponse{
				Error: 1,
				Msg:   err.Error(),
			}, nil
		}
		entries = append(entries, &pb.EntryOptionWithTime{
			CronTime: val.CronTime,
			Url:      val.Url,
			Headers:  headers,
			Body:     body,
			NextDate: val.NextDate.Unix(),
			LastDate: val.LastDate.Unix(),
		})
	}
	return &pb.GetResponse{
		Error: 0,
		Data: &pb.Information{
			Identity: opt.Identity,
			Start:    opt.Start,
			TimeZone: opt.TimeZone,
			Entries:  entries,
		},
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
