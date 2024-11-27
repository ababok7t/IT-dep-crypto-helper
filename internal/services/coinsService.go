package services

import (
	"crypto-helper/internal/domain"
	"errors"
	"fmt"
	"math"
	"strconv"
)

func (s *Service) UpdateCoinsInfo() error {
	coins, getError := s.Client.GetCoinsInfo()
	if getError != nil {
		return getError
	}

	s.CoinsCache.SetCoins(coins)

	return nil
}

func (s *Service) GetCoinInfo(symbol string) (domain.Coin, error) {
	coins, isFound := s.CoinsCache.GetCoins()

	if !isFound {
		return domain.Coin{}, errors.New("coin not found")
	}

	return coins[symbol], nil
}

func (s *Service) GetCoinForecast(symbol string) (string, error) {
	coin, err := s.GetCoinInfo(symbol)

	if err != nil {
		return "", err
	}

	priceNow, _ := strconv.ParseFloat(coin.PriceUsd, 64)

	change1H, _ := strconv.ParseFloat(coin.PercentChange1H, 64)
	change24H, _ := strconv.ParseFloat(coin.PercentChange24H, 64)
	change7D, _ := strconv.ParseFloat(coin.PercentChange7D, 64)

	price1H := priceNow / (1 + change1H/100)
	price24H := priceNow / (1 + change24H/100)
	price7D := priceNow / (1 + change7D/100)

	logProfit1 := math.Log(priceNow / price1H)
	logProfit2 := math.Log(price24H / price1H)
	logProfit3 := math.Log(price7D / price24H)

	midProfit := (logProfit1 + logProfit2 + logProfit3) / 3

	vol1 := math.Pow(logProfit1-midProfit, 2)
	vol2 := math.Pow(logProfit2-midProfit, 2)
	vol3 := math.Pow(logProfit3-midProfit, 2)

	volatility := math.Sqrt((vol1 + vol2 + vol3) / 3)

	coinForecast := priceNow*math.Exp(midProfit-1/2*math.Pow(volatility, 2)) + volatility*0.0000001

	return fmt.Sprintf("%.10f", coinForecast), nil
}
