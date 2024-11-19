package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)

	cache := Cache{
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		items:             items,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {

	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Created:    time.Now(),
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {

	c.RLock()

	defer c.RUnlock()

	item, isFound := c.items[key]

	if !isFound {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	_, isFound := c.items[key]

	if !isFound {
		return errors.New("key not found")
	}

	delete(c.items, key)

	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		keys := c.expiredKeys()

		if len(keys) > 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for key, value := range c.items {

		if time.Now().UnixNano() > value.Expiration && value.Expiration > 0 {
			keys = append(keys, key)
		}

	}

	return keys
}

func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, key := range keys {
		delete(c.items, key)
	}
}
