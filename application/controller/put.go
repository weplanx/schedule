package controller

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	pb "schedule-microservice/api"
)

func (c *controller) Put(_ context.Context, option *pb.Option) (*empty.Empty, error) {
	//var err error
	//entryOption := make(map[string]*types.EntryOption)
	//for taskID, val := range param.Entries {
	//	var headers map[string]string
	//	err = jsoniter.Unmarshal(val.Headers, &headers)
	//	if err != nil {
	//		return c.response(err)
	//	}
	//	var body interface{}
	//	err = jsoniter.Unmarshal(val.Body, &body)
	//	if err != nil {
	//		return c.response(err)
	//	}
	//	entryOption[taskID] = &types.EntryOption{
	//		CronTime: val.CronTime,
	//		Url:      val.Url,
	//		Headers:  headers,
	//		Body:     body,
	//		NextDate: time.Time{},
	//		LastDate: time.Time{},
	//	}
	//}
	//err = c.manager.Put(types.JobOption{
	//	Identity: param.Identity,
	//	TimeZone: param.TimeZone,
	//	Start:    param.Start,
	//	Entries:  entryOption,
	//})
	//if err != nil {
	//	return c.response(err)
	//}
	return nil, nil
}
