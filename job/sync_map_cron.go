package job

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type syncMapCron struct {
	sync.RWMutex
	Map map[string]*cron.Cron
}

func newSyncMapCron() *syncMapCron {
	c := new(syncMapCron)
	c.Map = make(map[string]*cron.Cron)
	return c
}

func (c *syncMapCron) Empty(identity string) bool {
	return c.Map[identity] == nil
}

func (c *syncMapCron) Get(identity string) *cron.Cron {
	c.RLock()
	data := c.Map[identity]
	c.RUnlock()
	return data
}

func (c *syncMapCron) Set(identity string, cron *cron.Cron) {
	c.Lock()
	c.Map[identity] = cron
	c.Unlock()
}

func (c *syncMapCron) Clear(identity string) {
	delete(c.Map, identity)
}
