package controller

import (
	"context"
	"encoding/json"
	"schedule-microservice/common"
	pb "schedule-microservice/router"
	"schedule-microservice/task"
)

type controller struct {
	tk *task.Task
	pb.UnimplementedRouterServer
}

func New(tk *task.Task) *controller {
	c := new(controller)
	c.tk = tk
	return c
}

func (c *controller) Get(ctx context.Context, req *pb.GetParameter) (*pb.GetResponse, error) {
	var err error
	opt := c.tk.Get(req.Identity)
	entries := make(map[string]*pb.EntryOptionWithTime)
	for key, val := range opt.Entries {
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
		entries[key] = &pb.EntryOptionWithTime{
			CronTime: val.CronTime,
			Url:      val.Url,
			Headers:  headers,
			Body:     body,
			NextDate: val.NextDate.Unix(),
			LastDate: val.LastDate.Unix(),
		}
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
	var err error
	var lists []*pb.Information
	for _, identity := range req.Identity {
		opt := c.tk.Get(identity)
		entries := make(map[string]*pb.EntryOptionWithTime)
		for key, val := range opt.Entries {
			var headers []byte
			headers, err = json.Marshal(val.Headers)
			if err != nil {
				return &pb.ListsResponse{
					Error: 1,
					Msg:   err.Error(),
				}, nil
			}
			var body []byte
			body, err = json.Marshal(val.Body)
			if err != nil {
				return &pb.ListsResponse{
					Error: 1,
					Msg:   err.Error(),
				}, nil
			}
			entries[key] = &pb.EntryOptionWithTime{
				CronTime: val.CronTime,
				Url:      val.Url,
				Headers:  headers,
				Body:     body,
				NextDate: val.NextDate.Unix(),
				LastDate: val.LastDate.Unix(),
			}
		}
		lists = append(lists, &pb.Information{
			Identity: opt.Identity,
			Start:    opt.Start,
			TimeZone: opt.TimeZone,
			Entries:  entries,
		})
	}
	return &pb.ListsResponse{
		Error: 0,
		Data:  lists,
	}, nil
}

func (c *controller) All(ctx context.Context, req *pb.NoParameter) (*pb.AllResponse, error) {
	return &pb.AllResponse{
		Error: 0,
		Data:  c.tk.All(),
	}, nil
}

func (c *controller) Put(ctx context.Context, req *pb.PutParameter) (*pb.Response, error) {
	var err error
	entries := make(map[string]*common.EntryOption)
	for key, val := range req.Entries {
		var headers map[string]string
		err = json.Unmarshal(val.Headers, headers)
		if err != nil {
			return &pb.Response{
				Error: 1,
				Msg:   err.Error(),
			}, nil
		}
		var body interface{}
		err = json.Unmarshal(val.Body, body)
		if err != nil {
			return &pb.Response{
				Error: 1,
				Msg:   err.Error(),
			}, nil
		}
		entries[key] = &common.EntryOption{
			CronTime: val.CronTime,
			Url:      val.Url,
			Headers:  headers,
			Body:     body,
			NextDate: nil,
			LastDate: nil,
		}
	}
	err = c.tk.Put(common.TaskOption{
		Identity: req.Identity,
		TimeZone: req.TimeZone,
		Start:    req.Start,
		Entries:  entries,
	})
	if err != nil {
		return &pb.Response{
			Error: 1,
			Msg:   err.Error(),
		}, nil
	}
	return &pb.Response{
		Error: 0,
		Msg:   "ok",
	}, nil
}

func (c *controller) Delete(ctx context.Context, req *pb.DeleteParameter) (*pb.Response, error) {
	err := c.tk.Delete(req.Identity)
	if err != nil {
		return &pb.Response{
			Error: 1,
			Msg:   err.Error(),
		}, nil
	}
	return &pb.Response{
		Error: 0,
		Msg:   "ok",
	}, nil
}

func (c *controller) Running(ctx context.Context, req *pb.RunningParameter) (*pb.Response, error) {
	err := c.tk.Running(req.Identity, req.Running)
	if err != nil {
		return &pb.Response{
			Error: 1,
			Msg:   err.Error(),
		}, nil
	}
	return &pb.Response{
		Error: 0,
		Msg:   "ok",
	}, nil
}
