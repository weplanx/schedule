package utils

import (
	"schedule-microservice/config/options"
	"sync"
)

type JobOptionMap struct {
	sync.RWMutex
	hashMap map[string]*options.JobOption
}

func NewJobOptionMap() *JobOptionMap {
	c := new(JobOptionMap)
	c.hashMap = make(map[string]*options.JobOption)
	return c
}

func (c *JobOptionMap) Put(identity string, option *options.JobOption) {
	c.Lock()
	c.hashMap[identity] = option
	c.Unlock()
}

func (c *JobOptionMap) Start(identity string, value bool) {
	c.Lock()
	c.hashMap[identity].Start = value
	c.Unlock()
}

func (c *JobOptionMap) Empty(identity string) bool {
	return c.hashMap[identity] == nil
}

func (c *JobOptionMap) Get(identity string) *options.JobOption {
	c.RLock()
	value := c.hashMap[identity]
	c.RUnlock()
	return value
}

func (c *JobOptionMap) Lists() map[string]*options.JobOption {
	return c.hashMap
}

func (c *JobOptionMap) Remove(identity string) {
	delete(c.hashMap, identity)
}
