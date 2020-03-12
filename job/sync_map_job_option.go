package job

import (
	"schedule-microservice/common"
	"sync"
)

type syncMapJobOption struct {
	sync.RWMutex
	Map map[string]*common.JobOption
}

func newSyncMapJobOption() *syncMapJobOption {
	c := new(syncMapJobOption)
	c.Map = make(map[string]*common.JobOption)
	return c
}

func (c *syncMapJobOption) Empty(identity string) bool {
	return c.Map[identity] == nil
}

func (c *syncMapJobOption) Get(identity string) *common.JobOption {
	c.RLock()
	data := c.Map[identity]
	c.RUnlock()
	return data
}

func (c *syncMapJobOption) Set(identity string, option *common.JobOption) {
	c.Lock()
	c.Map[identity] = option
	c.Unlock()
}

func (c *syncMapJobOption) Clear(identity string) {
	delete(c.Map, identity)
}
