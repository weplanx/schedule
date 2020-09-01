package utils

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type SyncMapEntryID struct {
	sync.RWMutex
	Map map[string]map[string]cron.EntryID
}

func NewSyncMapEntryID() *SyncMapEntryID {
	c := new(SyncMapEntryID)
	c.Map = make(map[string]map[string]cron.EntryID)
	return c
}

func (c *SyncMapEntryID) Empty(identity string) bool {
	return c.Map[identity] == nil
}

func (c *SyncMapEntryID) Set(identity string, taskID string, EntryID cron.EntryID) {
	c.Lock()
	if c.Map[identity] == nil {
		c.Map[identity] = make(map[string]cron.EntryID)
	}
	c.Map[identity][taskID] = EntryID
	c.Unlock()
}

func (c *SyncMapEntryID) Clear(identity string) {
	delete(c.Map, identity)
}
