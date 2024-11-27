package cache

import (
	"crypto-helper/internal/domain"
	"sync"
)

type CoinsCache struct {
	sync.RWMutex
	coins map[string]domain.Coin
}

func NewCoinsCache() *CoinsCache {

	return &CoinsCache{
		coins: make(map[string]domain.Coin),
	}

}

func (c *CoinsCache) SetCoins(value map[string]domain.Coin) {

	c.Lock()
	defer c.Unlock()

	c.coins = value
}

func (c *CoinsCache) GetCoins() (map[string]domain.Coin, bool) {

	c.RLock()
	defer c.RUnlock()

	if c.coins == nil {
		return nil, false
	}

	return c.coins, true
}
