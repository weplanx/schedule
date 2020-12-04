package utils

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type EntryIDMap struct {
	sync.RWMutex
	hashMap map[string]map[string]cron.EntryID
}

func NewEntryIDMap() *EntryIDMap {
	c := new(EntryIDMap)
	c.hashMap = make(map[string]map[string]cron.EntryID)
	return c
}

func (c *EntryIDMap) Put(identity string, taskID string, EntryID cron.EntryID) {
	c.Lock()
	if c.hashMap[identity] == nil {
		c.hashMap[identity] = make(map[string]cron.EntryID)
	}
	c.hashMap[identity][taskID] = EntryID
	c.Unlock()
}

func (c *EntryIDMap) Empty(identity string) bool {
	return c.hashMap[identity] == nil
}

func (c *EntryIDMap) Get(identity string) map[string]cron.EntryID {
	c.RLock()
	value := c.hashMap[identity]
	c.RUnlock()
	return value
}

func (c *EntryIDMap) Lists() map[string]map[string]cron.EntryID {
	return c.hashMap
}

func (c *EntryIDMap) Remove(identity string) {
	delete(c.hashMap, identity)
}
