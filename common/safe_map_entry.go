package common

import (
	"sync"
)

type safeMapEntry struct {
	sync.RWMutex
	Map map[string]*EntryOption
}

func NewSafeMapEntry() *safeMapEntry {
	safe := new(safeMapEntry)
	safe.Map = make(map[string]*EntryOption)
	return safe
}

func (s *safeMapEntry) Clear(identity string) {
	delete(s.Map, identity)
}

func (s *safeMapEntry) Get(identity string) *EntryOption {
	s.RLock()
	entry := s.Map[identity]
	s.RUnlock()
	return entry
}

func (s *safeMapEntry) Set(identity string, entry *EntryOption) {
	s.Lock()
	s.Map[identity] = entry
	s.Unlock()
}
