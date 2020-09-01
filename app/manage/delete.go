package manage

func (c *JobsManager) Delete(identity string) (err error) {
	c.termination(identity)
	c.options.Clear(identity)
	c.runtime.Clear(identity)
	c.entryIDSet.Clear(identity)
	return
}
