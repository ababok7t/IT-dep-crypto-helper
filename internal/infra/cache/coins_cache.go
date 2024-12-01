package cache

import (
	"crypto-helper/internal/domain"
	"sync"
)

type CoinsCache struct {
	mtx          sync.RWMutex //mtx// //расширяем чтобы запихнуть//
	coins        map[string]domain.Coin
	coinsSymbols []string
}

func NewCoinsCache() *CoinsCache {
	return &CoinsCache{
		coins: make(map[string]domain.Coin),
	}
}

func (c *CoinsCache) SetCoins(coins map[string]domain.Coin, coinsSymbols []string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.coins = coins
	c.coinsSymbols = coinsSymbols
}

func (c *CoinsCache) GetCoins() (map[string]domain.Coin, []string, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if c.coins == nil {
		return nil, nil, false
	}

	return c.coins, c.coinsSymbols, true
}
