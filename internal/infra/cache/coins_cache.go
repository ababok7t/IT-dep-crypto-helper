package cache

import (
	"crypto-helper/internal/domain"
	"sync"
)

type CoinsCache struct {
	sync.RWMutex
	coins        map[string]domain.Coin
	coinsSymbols []string
}

func NewCoinsCache() *CoinsCache {

	return &CoinsCache{
		coins: make(map[string]domain.Coin),
	}

}

func (c *CoinsCache) SetCoins(coins map[string]domain.Coin, coinsSymbols []string) {

	c.Lock()
	defer c.Unlock()

	c.coins = coins
	c.coinsSymbols = coinsSymbols
}

func (c *CoinsCache) GetCoins() (map[string]domain.Coin, []string, bool) {

	c.RLock()
	defer c.RUnlock()

	if c.coins == nil {
		return nil, nil, false
	}

	return c.coins, c.coinsSymbols, true
}
