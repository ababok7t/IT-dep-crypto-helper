package cache

import (
	"sync"
)

type StatesCache struct {
	mtx    sync.RWMutex
	states map[string]string
}

func NewStatesCache() *StatesCache {
	return &StatesCache{
		states: make(map[string]string),
	}
}

func (c *StatesCache) SetState(id string, state string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.states[id] = state
}

func (c *StatesCache) GetState(id string) string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.states[id]
}
