package controller

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	jsoniter "github.com/json-iterator/go"
	pb "schedule-microservice/api"
	"schedule-microservice/config/options"
)

func (c *controller) Put(_ context.Context, option *pb.Option) (_ *empty.Empty, err error) {
	opt := options.JobOption{
		Identity: option.Id,
		TimeZone: option.TimeZone,
		Start:    option.Start,
		Entries:  nil,
	}
	opt.Entries = make(map[string]*options.EntryOption)
	for key, value := range option.Entries {
		if err = validator.New().Var(string(value.Headers), "json"); err != nil {
			return
		}
		var headers map[string]string
		jsoniter.Unmarshal(value.Headers, &headers)
		if err = validator.New().Var(string(value.Body), "json"); err != nil {
			return
		}
		var body map[string]interface{}
		jsoniter.Unmarshal(value.Body, &body)
		opt.Entries[key] = &options.EntryOption{
			CronTime: value.CronTime,
			Url:      value.Url,
			Headers:  headers,
			Body:     body,
		}
	}
	if err = c.Job.Put(opt); err != nil {
		return
	}
	return &empty.Empty{}, nil
}
