package common

import "time"

type Job struct {
	Mode   string `msgpack:"mode"`
	Spec   string `msgpack:"spec"`
	Option Option `msgpack:"option"`

	State `msgpack:"state"`
}

type State struct {
	Next time.Time `msgpack:"next"`
	Prev time.Time `msgpack:"prev"`
}

type Option interface{}

type HttpOption struct {
	Url     string                 `msgpack:"url"`
	Headers map[string]string      `msgpack:"headers"`
	Body    map[string]interface{} `msgpack:"body"`
}

// HttpJob HTTP回调
func HttpJob(spec string, option HttpOption) Job {
	return Job{
		Mode:   "HTTP",
		Spec:   spec,
		Option: option,
	}
}
