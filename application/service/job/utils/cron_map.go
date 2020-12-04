package utils

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type CronMap struct {
	sync.RWMutex
	hashMap map[string]*cron.Cron
}

func NewCronMap() *CronMap {
	c := new(CronMap)
	c.hashMap = make(map[string]*cron.Cron)
	return c
}

func (c *CronMap) Put(identity string, cron *cron.Cron) {
	c.Lock()
	c.hashMap[identity] = cron
	c.Unlock()
}

func (c *CronMap) Empty(identity string) bool {
	return c.hashMap[identity] == nil
}

func (c *CronMap) Get(identity string) *cron.Cron {
	c.RLock()
	value := c.hashMap[identity]
	c.RUnlock()
	return value
}

func (c *CronMap) Lists() map[string]*cron.Cron {
	return c.hashMap
}

func (c *CronMap) Remove(identity string) {
	delete(c.hashMap, identity)
}
