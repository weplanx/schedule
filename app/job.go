package app

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
func HttpJob(spec string, option map[string]interface{}) *Job {
	return NewJob(spec, func() {

	})
}
