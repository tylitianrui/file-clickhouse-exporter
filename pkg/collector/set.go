package collector

import "sync"

const InitSetSize = 1 << 4

type Set struct {
	vals map[interface{}]bool
	mu   sync.Mutex
}

func NewSet() *Set {
	return &Set{
		vals: make(map[interface{}]bool, InitSetSize),
	}
}

func (s *Set) Add(item interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vals[item] = true
}

func (s *Set) AllItems() (items []interface{}) {
	for item, flag := range s.vals {
		if flag {
			items = append(items, item)
		}
	}
	return
}
