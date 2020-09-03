package manage

func (c *JobsManager) Delete(identity string) (err error) {
	if c.options.Empty(identity) {
		return
	}
	c.termination(identity)
	c.options.Clear(identity)
	c.runtime.Clear(identity)
	c.entryIDSet.Clear(identity)
	return c.schema.Delete(identity)
}
