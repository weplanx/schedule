package common

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/weplanx/schedule/model"
	"go.uber.org/zap"
)

type Job struct {
	*Inject

	spec string
	cmd  func()
}

// NewJob 创建工作
func NewJob(i *Inject, spec string) *Job {
	return &Job{
		Inject: i,
		spec:   spec,
	}
}

// HttpWorker Http 回调工作
func (x *Job) HttpWorker(ctx context.Context, option model.HttpJob) *Job {
	x.cmd = func() {
		client := resty.New()
		resp, err := client.R().
			SetHeaders(option.Headers).
			SetBody(option.Body).
			Post(option.Url)
		if err != nil {
			x.Log.Error("发起失败",
				zap.String("url", option.Url),
				zap.Error(err),
				zap.String("time", resp.Time().String()),
			)
			return
		}
		x.Log.Debug("发起成功",
			zap.String("url", option.Url),
			zap.Int("status", resp.StatusCode()),
			zap.ByteString("body", resp.Body()),
			zap.String("time", resp.Time().String()),
		)
		if err = x.Transfer.Publish(ctx, x.Values.Transfer.Topic, map[string]interface{}{
			"key":    "schedule",
			"node":   x.Values.Node,
			"url":    option.Url,
			"status": resp.StatusCode(),
			"body":   resp.Body(),
			"time":   resp.Time(),
		}); err != nil {
			return
		}
	}
	return x
}
