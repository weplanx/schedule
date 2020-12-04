package job

func (c *Job) Running(identity string, status bool) error {
	if c.Options.Empty(identity) {
		return NotExists
	}
	c.Options.Start(identity, status)
	runtime := c.Runtime.Get(identity)
	if status {
		runtime.Start()
	} else {
		runtime.Stop()
	}
	return c.Schema.Update(*c.Options.Get(identity))
}
