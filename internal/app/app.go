package app

import (
	"crypto-helper/internal/infra/cache"
	"crypto-helper/internal/services"
	"time"
)

func RunApp() {
	cache := cache.New(time.Second*30, time.Minute)

	coins, gettingError := services.GetCoinsInfo()

	if gettingError != nil {
		panic(gettingError)
	}

	for _, coin := range coins {
		cache.Set(coin.Symbol, coin, time.Second*30)
	}

	return
}
