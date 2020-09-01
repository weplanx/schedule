package utils

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type SyncMapCron struct {
	sync.RWMutex
	Map map[string]*cron.Cron
}

func NewSyncMapCron() *SyncMapCron {
	c := new(SyncMapCron)
	c.Map = make(map[string]*cron.Cron)
	return c
}

func (c *SyncMapCron) Empty(identity string) bool {
	return c.Map[identity] == nil
}

func (c *SyncMapCron) Get(identity string) *cron.Cron {
	c.RLock()
	data := c.Map[identity]
	c.RUnlock()
	return data
}

func (c *SyncMapCron) Set(identity string, cron *cron.Cron) {
	c.Lock()
	c.Map[identity] = cron
	c.Unlock()
}

func (c *SyncMapCron) Clear(identity string) {
	delete(c.Map, identity)
}
