package utils

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type SyncMapEntryIDSet struct {
	sync.RWMutex
	Map map[string]map[string]cron.EntryID
}

func NewSyncMapEntryIDSet() *SyncMapEntryIDSet {
	c := new(SyncMapEntryIDSet)
	c.Map = make(map[string]map[string]cron.EntryID)
	return c
}

func (c *SyncMapEntryIDSet) Empty(identity string) bool {
	return c.Map[identity] == nil
}

func (c *SyncMapEntryIDSet) Set(identity string, task string, EntryID cron.EntryID) {
	c.Lock()
	if c.Map[identity] == nil {
		c.Map[identity] = make(map[string]cron.EntryID)
	}
	c.Map[identity][task] = EntryID
	c.Unlock()
}

func (c *SyncMapEntryIDSet) Clear(identity string) {
	delete(c.Map, identity)
}
