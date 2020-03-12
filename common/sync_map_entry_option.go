package common

import (
	"sync"
)

type syncMapEntryOption struct {
	sync.RWMutex
	Map map[string]*EntryOption
}

func NewSyncMapEntryOption() *syncMapEntryOption {
	c := new(syncMapEntryOption)
	c.Map = make(map[string]*EntryOption)
	return c
}

func (c *syncMapEntryOption) Get(identity string) *EntryOption {
	c.RLock()
	data := c.Map[identity]
	c.RUnlock()
	return data
}

func (c *syncMapEntryOption) Set(identity string, entry *EntryOption) {
	c.Lock()
	c.Map[identity] = entry
	c.Unlock()
}

func (c *syncMapEntryOption) Clear(identity string) {
	delete(c.Map, identity)
}
