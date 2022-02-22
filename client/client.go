package client

import (
	"context"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weplanx/schedule/api"
	"github.com/weplanx/schedule/model"
	"google.golang.org/grpc"
)

type Schedule struct {
	client api.APIClient
	conn   *grpc.ClientConn
}

func New(addr string, opts ...grpc.DialOption) (x *Schedule, err error) {
	x = new(Schedule)
	if x.conn, err = grpc.Dial(addr, opts...); err != nil {
		return
	}
	x.client = api.NewAPIClient(x.conn)
	return
}

func (x *Schedule) Close() error {
	return x.conn.Close()
}

// HttpWorker HTTP回调
// @spec Crontab 规则
// @option HTTP 配置
func HttpWorker(spec string, option model.HttpJob) (job *api.Job, err error) {
	job = &api.Job{
		Spec: spec,
		Mode: "HTTP",
	}
	if job.Option, err = msgpack.Marshal(option); err != nil {
		return
	}
	return
}

// Put 更新任务
// @key 唯一标识
// @jobs 工作流
func (x *Schedule) Put(ctx context.Context, key string, jobs ...*api.Job) (err error) {
	if _, err = x.client.Put(ctx, &api.Schedule{
		Key:  key,
		Jobs: jobs,
	}); err != nil {
		return
	}
	return
}

// Get 获取任务
// @keys 唯一标识数组
func (x *Schedule) Get(ctx context.Context, keys []string) (data map[string]*api.Schedule, err error) {
	var rep *api.GetReply
	if rep, err = x.client.Get(ctx, &api.GetRequest{Keys: keys}); err != nil {
		return
	}
	data = rep.GetData()
	return
}

// Delete 删除任务
// @key 唯一标识
func (x *Schedule) Delete(ctx context.Context, key string) (err error) {
	if _, err = x.client.Delete(ctx, &api.DeleteRequest{Key: key}); err != nil {
		return
	}
	return
}
