package app

import (
	"github.com/go-resty/resty/v2"
	"github.com/weplanx/schedule/model"
	"go.uber.org/zap"
)

type Job struct {
	spec string
	cmd  func()
}

// NewJob 创建工作
func NewJob(spec string, cmd func()) *Job {
	return &Job{
		spec: spec,
		cmd:  cmd,
	}
}

// HttpJob Http回调工作
func HttpJob(spec string, option model.HttpJob, log *zap.Logger) *Job {
	return NewJob(spec, func() {
		client := resty.New()
		resp, err := client.R().
			SetHeaders(option.Headers).
			SetBody(option.Body).
			Post(option.Url)
		if err != nil {
			log.Error("发起失败",
				zap.String("url", option.Url),
				zap.Int("status", resp.StatusCode()),
				zap.Error(err),
				zap.String("time", resp.Time().String()),
			)
			return
		}
		log.Debug("发起成功",
			zap.String("url", option.Url),
			zap.Int("status", resp.StatusCode()),
			zap.ByteString("body", resp.Body()),
			zap.String("time", resp.Time().String()),
		)
	})
}
