package job

func (c *Job) Delete(identity string) (err error) {
	if c.Options.Empty(identity) {
		return
	}
	c.termination(identity)
	c.Options.Remove(identity)
	c.Runtime.Remove(identity)
	c.EntryIDSet.Remove(identity)
	return c.Schema.Delete(identity)
}
