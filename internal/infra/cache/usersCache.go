package cache

import (
	"crypto-helper/internal/domain"
	"sync"
	"time"
)

type UsersCache struct {
	sync.RWMutex
	users map[string]domain.User
}

func NewUsersCache() *UsersCache {

	return &UsersCache{
		users: make(map[string]domain.User),
	}

}

func (c *UsersCache) SetUser(user domain.User) {

	c.Lock()
	defer c.Unlock()

	c.users[user.UserId] = user
}

func (c *UsersCache) GetUser(userId string) (domain.User, bool) {

	c.RLock()
	defer c.RUnlock()

	user, isFound := c.users[userId]
	return user, isFound
}

func (c *UsersCache) GetAllUsers() map[string]domain.User {

	c.RLock()
	defer c.RUnlock()

	return c.users
}

func (c *UsersCache) SetPriceAlert(userId string, alert domain.PriceAlert, ttl time.Duration) {

	user, isFound := c.GetUser(userId)

	if !isFound {
		user = domain.User{UserId: userId, Collection: []domain.Coin{}}
	}

	user.PriceAlertsList = append(user.PriceAlertsList, alert)

	c.SetUser(user)

	go func() {
		time.Sleep(ttl)
		c.DeletePriceAlert(userId, alert.Symbol)
	}()
}

func (c *UsersCache) GetPriceAlert(userId string, coinSymbol string) (domain.PriceAlert, bool) {
	user, isFound := c.GetUser(userId)

	if !isFound {
		return domain.PriceAlert{}, false
	}

	for _, priceAlert := range user.PriceAlertsList {
		if priceAlert.Symbol == coinSymbol {
			return priceAlert, true
		}
	}

	return domain.PriceAlert{}, false
}

func (c *UsersCache) DeletePriceAlert(userId string, coinSymbol string) {

	user, isFound := c.GetUser(userId)

	if !isFound {
		return
	}

	var updatedPriceAlerts []domain.PriceAlert
	for _, priceAlert := range user.PriceAlertsList {
		if coinSymbol != priceAlert.Symbol {
			updatedPriceAlerts = append(updatedPriceAlerts, priceAlert)
		}
	}

	user.PriceAlertsList = updatedPriceAlerts
	c.SetUser(user)
}

func (c *UsersCache) SetCollectionItem(userId string, coin domain.Coin) {

	user, isFound := c.GetUser(userId)
	if !isFound {
		user = domain.User{UserId: userId, PriceAlertsList: []domain.PriceAlert{}}
	}

	user.Collection = append(user.Collection, coin)

	c.SetUser(user)
}

func (c *UsersCache) DeleteCollectionItem(userId string, coin domain.Coin) {
	user, isFound := c.GetUser(userId)

	if !isFound {
		return
	}

	var updatedCollectionItems []domain.Coin
	for _, item := range user.Collection {
		if item.Symbol != coin.Symbol {
			updatedCollectionItems = append(updatedCollectionItems, item)
		}
	}

	user.Collection = updatedCollectionItems
	c.SetUser(user)
}

func (c *UsersCache) GetAllCollectionItems(userId string) []domain.Coin {
	user, isFound := c.GetUser(userId)

	if !isFound {
		user = domain.User{UserId: userId, Collection: []domain.Coin{}, PriceAlertsList: []domain.PriceAlert{}}
		c.SetUser(user)
		return user.Collection
	}

	return user.Collection
}
