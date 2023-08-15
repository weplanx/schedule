package app

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/collector/transfer"
	"github.com/weplanx/workflow/typ"
	"github.com/weplanx/workflow/worker/common"
	"go.uber.org/zap"
	"time"
)

type App struct {
	*common.Inject
}

type M = map[string]interface{}

func Initialize(i *common.Inject) *App {
	return &App{
		Inject: i,
	}
}

func (x *App) Run(ctx context.Context) (err error) {
	if err = x.Transfer.Set(ctx, transfer.StreamOption{
		Key: "logset_jobs",
	}); err != nil {
		return
	}
	subj := fmt.Sprintf(`%s.jobs.*`, x.V.Namespace)
	queue := fmt.Sprintf(`%s:worker`, x.V.Namespace)
	if _, err = x.Nats.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		var job typ.Job
		if err = msgpack.Unmarshal(msg.Data, &job); err != nil {
			return
		}
		switch job.Mode {
		case "HTTP":
			x.HTTPMode(job)
			break
		}
	}); err != nil {
		return
	}
	return
}

func (x *App) HTTPMode(job typ.Job) (err error) {
	httpClient := resty.New()
	httpClient.JSONMarshal = sonic.Marshal
	httpClient.JSONUnmarshal = sonic.Unmarshal
	var option typ.HttpOption
	if err = mapstructure.Decode(job.Option, &option); err != nil {
		x.Log.Error("HttpOption:fail",
			zap.Any("option", job.Option),
			zap.Error(err),
		)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := httpClient.R().
		SetContext(ctx).
		SetHeaders(option.Headers).
		SetBody(option.Body).
		Post(option.Url)
	now := time.Now()
	x.Transfer.Publish(ctx, "logset_jobs", transfer.Payload{
		Timestamp: now,
		Data: M{
			"metadata": M{
				"key":   job.Key,
				"index": job.Index,
				"mode":  job.Mode,
				"url":   option.Url,
			},
			"headers": option.Headers,
			"body":    option.Body,
			"response": M{
				"status": resp.StatusCode(),
				"body":   resp.Result(),
			},
		},
		XData: M{},
	})
	if err != nil {
		x.Log.Error("Http:fail",
			zap.String("key", job.Key),
			zap.Int("index", job.Index),
			zap.Any("headers", option.Headers),
			zap.Any("body", option.Body),
			zap.Error(err),
		)
		return
	}
	x.Log.Debug("Http:ok",
		zap.String("key", job.Key),
		zap.Int("index", job.Index),
		zap.Any("headers", option.Headers),
		zap.Any("body", option.Body),
		zap.Any("response", M{
			"status": resp.StatusCode(),
			"body":   resp,
		}),
	)
	return
}
