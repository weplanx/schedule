package utiliy

type Job struct {
	Mode string `msgpack:"mode"`
	Rule string `msgpack:"rule"`
	Spec Spec   `msgpack:"spec"`
}

type Spec interface{}

type HttpSpec struct {
	Url     string                 `msgpack:"url"`
	Headers map[string]string      `msgpack:"headers"`
	Body    map[string]interface{} `msgpack:"body"`
}

// HttpJob HTTP回调
func HttpJob(rule string, spec HttpSpec) Job {
	return Job{
		Mode: "HTTP",
		Rule: rule,
		Spec: spec,
	}
}
