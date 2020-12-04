package controller

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	pb "schedule-microservice/api"
	"schedule-microservice/application/common"
	"schedule-microservice/config/options"
)

type controller struct {
	pb.UnimplementedAPIServer
	*common.Dependency
}

func New(dep *common.Dependency) *controller {
	c := new(controller)
	c.Dependency = dep
	return c
}

func (c *controller) find(identity string) (data *pb.Option, err error) {
	var option *options.JobOption
	var runtime *cron.Cron
	var entryID map[string]cron.EntryID
	if option, runtime, entryID, err = c.Job.Get(identity); err != nil {
		return
	}
	data = &pb.Option{
		Id:       option.Identity,
		TimeZone: option.TimeZone,
		Start:    option.Start,
		Entries:  nil,
	}
	data.Entries = make(map[string]*pb.Entry)
	for key, value := range entryID {
		entry := runtime.Entry(value)
		entryOption := option.Entries[key]
		var headers []byte
		jsoniter.Marshal(entryOption.Headers)
		var body []byte
		jsoniter.Marshal(entryOption.Body)
		data.Entries[key] = &pb.Entry{
			CronTime: entryOption.CronTime,
			Url:      entryOption.Url,
			Headers:  headers,
			Body:     body,
			NextDate: entry.Next.Unix(),
			LastDate: entry.Prev.Unix(),
		}
	}
	return data, nil
}
