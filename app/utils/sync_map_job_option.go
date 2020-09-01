package utils

import (
	"schedule-microservice/app/types"
	"sync"
)

type SyncMapJobOption struct {
	sync.RWMutex
	Map map[string]*types.JobOption
}

func NewSyncMapJobOption() *SyncMapJobOption {
	c := new(SyncMapJobOption)
	c.Map = make(map[string]*types.JobOption)
	return c
}

func (c *SyncMapJobOption) Empty(identity string) bool {
	return c.Map[identity] == nil
}

func (c *SyncMapJobOption) Get(identity string) *types.JobOption {
	c.RLock()
	data := c.Map[identity]
	c.RUnlock()
	return data
}

func (c *SyncMapJobOption) Set(identity string, option *types.JobOption) {
	c.Lock()
	c.Map[identity] = option
	c.Unlock()
}

func (c *SyncMapJobOption) Clear(identity string) {
	delete(c.Map, identity)
}
