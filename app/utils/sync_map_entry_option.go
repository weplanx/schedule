package utils

import (
	"schedule-microservice/app/types"
	"sync"
)

type SyncMapEntryOption struct {
	sync.RWMutex
	Map map[string]*types.EntryOption
}

func NewSyncMapEntryOption() *SyncMapEntryOption {
	c := new(SyncMapEntryOption)
	c.Map = make(map[string]*types.EntryOption)
	return c
}

func (c *SyncMapEntryOption) Get(identity string) *types.EntryOption {
	c.RLock()
	data := c.Map[identity]
	c.RUnlock()
	return data
}

func (c *SyncMapEntryOption) Set(identity string, entry *types.EntryOption) {
	c.Lock()
	c.Map[identity] = entry
	c.Unlock()
}

func (c *SyncMapEntryOption) Clear(identity string) {
	delete(c.Map, identity)
}
