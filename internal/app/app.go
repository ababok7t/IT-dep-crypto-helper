package app

import (
	"crypto-helper/internal/infra/cache"
	"crypto-helper/internal/services"
	"fmt"
	"time"
)

func RunApp() {
	appCache := cache.New(time.Second*30, time.Second*45)
	err := services.UpdateCoinsInfo(appCache)
	if err != nil {
		return
	}

	fmt.Print(services.GetCoinForecast("BTC", appCache))
	fmt.Print(services.GetCoinInfo("BTC", appCache))
	fmt.Print(services.GetCoinForecast("ETH", appCache))
	fmt.Print(services.GetCoinInfo("ETH", appCache))

	<-time.After(time.Second * 7)

	fmt.Print(services.GetCoinForecast("ETH", appCache))

}
