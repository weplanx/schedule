package job

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"schedule-microservice/application/common/actions"
	"schedule-microservice/application/common/typ"
	"schedule-microservice/config/options"
	"time"
)

func (c *Job) Put(option options.JobOption) (err error) {
	identity := option.Identity
	var timezone *time.Location
	if timezone, err = time.LoadLocation(option.TimeZone); err != nil {
		return
	}
	c.termination(identity)
	c.Options.Put(identity, &option)
	job := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	c.Runtime.Put(identity, job)
	for taskID := range option.Entries {
		go c.addTask(identity, taskID)
	}
	if option.Start {
		job.Start()
	}
	return c.Schema.Update(option)
}

func (c *Job) addTask(identity string, taskID string) {
	var err error
	option := c.Options.Get(identity).Entries[taskID]
	var EntryID cron.EntryID
	if EntryID, err = c.Runtime.Get(identity).AddFunc(option.CronTime, func() {
		content := typ.Log{
			Identity: identity,
			Task:     taskID,
			Url:      option.Url,
			Header:   option.Headers,
			Body:     option.Body,
			Time:     time.Now().Unix(),
		}
		body, errs := actions.Fetch(option.Url, option.Headers, option.Body)
		if len(errs) != 0 {
			info := make([]string, len(errs))
			for index, value := range errs {
				info[index] = value.Error()
			}
			content.Status = false
			content.Response = map[string]interface{}{
				"errs": info,
			}
		} else {
			resBody := string(body)
			if err = validator.New().Var(resBody, "json"); err != nil {
				content.Response = map[string]interface{}{
					"raw": resBody,
				}
			} else {
				jsoniter.Unmarshal(body, &content.Response)
			}
			content.Status = true
		}
		go c.Transfer.Push(content)
		go func() {
			var logger *zap.Logger
			if logger, err = c.Filelog.NewLogger(identity); err != nil {
				return
			}
			fields := []zap.Field{
				zap.String("task", content.Task),
				zap.String("url", content.Url),
				zap.Any("header", content.Header),
				zap.Any("body", content.Body),
				zap.Any("response", content.Response),
			}
			if content.Status {
				logger.Info(identity, fields...)
			} else {
				logger.Error(identity, fields...)
			}
		}()
	}); err != nil {
		return
	}
	c.EntryIDSet.Put(identity, taskID, EntryID)
}
