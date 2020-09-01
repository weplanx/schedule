package manage

func (c *JobsManager) Running(identity string, running bool) (err error) {
	if err = c.empty(identity); err != nil {
		return
	}
	c.options.Map[identity].Start = running
	if running {
		c.runtime.Map[identity].Start()
	} else {
		c.runtime.Map[identity].Stop()
	}
	return
}
