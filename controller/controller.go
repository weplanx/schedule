package controller

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"schedule-microservice/common"
	"schedule-microservice/job"
	pb "schedule-microservice/router"
	"time"
)

type controller struct {
	job *job.Job
	pb.UnimplementedRouterServer
}

func New(job *job.Job) *controller {
	c := new(controller)
	c.job = job
	return c
}

func (c *controller) Get(ctx context.Context, req *pb.GetParameter) (*pb.GetResponse, error) {
	var err error
	opt := c.job.Get(req.Identity)
	if opt == nil {
		return &pb.GetResponse{
			Error: 0,
			Data:  nil,
		}, nil
	}
	entries := make(map[string]*pb.EntryOptionWithTime)
	for key, val := range opt.Entries {
		var headers []byte
		headers, err = jsoniter.Marshal(val.Headers)
		if err != nil {
			return &pb.GetResponse{
				Error: 1,
				Msg:   err.Error(),
			}, nil
		}
		var body []byte
		body, err = jsoniter.Marshal(val.Body)
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
	var entries map[string]*pb.EntryOptionWithTime
	for _, identity := range req.Identity {
		opt := c.job.Get(identity)
		if opt == nil {
			continue
		}
		entries = make(map[string]*pb.EntryOptionWithTime)
		for key, val := range opt.Entries {
			var headers []byte
			headers, err = jsoniter.Marshal(val.Headers)
			if err != nil {
				return &pb.ListsResponse{
					Error: 1,
					Msg:   err.Error(),
				}, nil
			}
			var body []byte
			body, err = jsoniter.Marshal(val.Body)
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
		Data:  c.job.All(),
	}, nil
}

func (c *controller) Put(ctx context.Context, req *pb.PutParameter) (*pb.Response, error) {
	var err error
	entryOption := common.NewSyncMapEntryOption()
	for key, val := range req.Entries {
		var headers map[string]string
		err = jsoniter.Unmarshal(val.Headers, &headers)
		if err != nil {
			return &pb.Response{
				Error: 1,
				Msg:   err.Error(),
			}, nil
		}
		var body interface{}
		err = jsoniter.Unmarshal(val.Body, &body)
		if err != nil {
			return &pb.Response{
				Error: 1,
				Msg:   err.Error(),
			}, nil
		}
		entryOption.Set(key, &common.EntryOption{
			CronTime: val.CronTime,
			Url:      val.Url,
			Headers:  headers,
			Body:     body,
			NextDate: time.Time{},
			LastDate: time.Time{},
		})
	}
	err = c.job.Put(common.JobOption{
		Identity: req.Identity,
		TimeZone: req.TimeZone,
		Start:    req.Start,
		Entries:  entryOption.Map,
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
	err := c.job.Delete(req.Identity)
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
	err := c.job.Running(req.Identity, req.Running)
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
