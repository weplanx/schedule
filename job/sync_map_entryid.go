package job

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type syncMapEntryID struct {
	sync.RWMutex
	Map map[string]map[string]cron.EntryID
}

func newSyncMapEntryID() *syncMapEntryID {
	c := new(syncMapEntryID)
	c.Map = make(map[string]map[string]cron.EntryID)
	return c
}

func (c *syncMapEntryID) Empty(identity string) bool {
	return c.Map[identity] == nil
}

func (c *syncMapEntryID) GetJobEntryID(identity string) map[string]cron.EntryID {
	c.RLock()
	data := c.Map[identity]
	c.RUnlock()
	return data
}

func (c *syncMapEntryID) GetTaskEntryID(identity string, task string) cron.EntryID {
	c.RLock()
	data := c.Map[identity][task]
	c.RUnlock()
	return data
}

func (c *syncMapEntryID) Set(identity string, task string, EntryID cron.EntryID) {
	c.Lock()
	if c.Map[identity] == nil {
		c.Map[identity] = make(map[string]cron.EntryID)
	}
	c.Map[identity][task] = EntryID
	c.Unlock()
}

func (c *syncMapEntryID) Clear(identity string) {
	delete(c.Map, identity)
}
