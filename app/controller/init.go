package controller

import (
	jsoniter "github.com/json-iterator/go"
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

func (c *controller) find(identity string) (information *pb.Information, err error) {
	option, err := c.manager.GetOption(identity)
	if err != nil {
		return
	}
	job, err := c.manager.GetJob(identity)
	if err != nil {
		return
	}
	entryIDs, err := c.manager.GetEntryIDSet(identity)
	if err != nil {
		return
	}
	entries := make(map[string]*pb.EntryOptionWithTime)
	for taskID, entryID := range entryIDs {
		task := job.Entry(entryID)
		entryOption := option.Entries[taskID]
		headers, err := jsoniter.Marshal(entryOption.Headers)
		if err != nil {
			return
		}
		body, err := jsoniter.Marshal(entryOption.Body)
		if err != nil {
			return
		}
		entries[taskID] = &pb.EntryOptionWithTime{
			CronTime: entryOption.CronTime,
			Url:      entryOption.Url,
			Headers:  headers,
			Body:     body,
			NextDate: task.Next.Unix(),
			LastDate: task.Prev.Unix(),
		}
	}
	information = &pb.Information{
		Identity: identity,
		Start:    option.Start,
		TimeZone: option.TimeZone,
		Entries:  entries,
	}
	return
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
