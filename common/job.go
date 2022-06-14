package common

import "time"

type Job struct {
	// 触发模式
	Mode string `msgpack:"mode"`

	// 时间规格
	Spec string `msgpack:"spec"`

	// 配置
	Option Option `msgpack:"option"`

	// 状态
	State `msgpack:"state"`
}

type Option interface{}

type HttpOption struct {
	// 网络回调地址
	Url string `msgpack:"url"`

	// 请求头部
	Headers map[string]string `msgpack:"headers"`

	// 请求体
	Body map[string]interface{} `msgpack:"body"`
}

// HttpJob HTTP回调
func HttpJob(spec string, option HttpOption) Job {
	return Job{
		Mode:   "HTTP",
		Spec:   spec,
		Option: option,
	}
}

type State struct {
	// 下次时间
	Next time.Time `msgpack:"next"`

	// 上次时间
	Prev time.Time `msgpack:"prev"`
}
