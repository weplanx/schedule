package controller

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"schedule-microservice/app/types"
	pb "schedule-microservice/router"
	"time"
)

func (c *controller) Put(ctx context.Context, param *pb.PutParameter) (*pb.Response, error) {
	var err error
	entryOption := make(map[string]*types.EntryOption)
	for taskID, val := range param.Entries {
		var headers map[string]string
		err = jsoniter.Unmarshal(val.Headers, &headers)
		if err != nil {
			return c.response(err)
		}
		var body interface{}
		err = jsoniter.Unmarshal(val.Body, &body)
		if err != nil {
			return c.response(err)
		}
		entryOption[taskID] = &types.EntryOption{
			CronTime: val.CronTime,
			Url:      val.Url,
			Headers:  headers,
			Body:     body,
			NextDate: time.Time{},
			LastDate: time.Time{},
		}
	}
	err = c.manager.Put(types.JobOption{
		Identity: param.Identity,
		TimeZone: param.TimeZone,
		Start:    param.Start,
		Entries:  entryOption,
	})
	if err != nil {
		return c.response(err)
	}
	return c.response(nil)
}
