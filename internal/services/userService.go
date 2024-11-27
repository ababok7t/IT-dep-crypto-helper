package services

import (
	"crypto-helper/internal/domain"
	"strconv"
	"time"
)

func (s *Service) AddPriceAlert(userId string, coinSymbol string, stopLimit string) {
	coin, err := s.GetCoinInfo(coinSymbol)
	if err != nil {
		return
	}

	price, priceParseErr := strconv.ParseFloat(coin.PriceUsd, 64)
	if priceParseErr != nil {
		return
	}

	stop, stopLimitParseErr := strconv.ParseFloat(stopLimit, 64)
	if stopLimitParseErr != nil {
		return
	}

	var status string

	if price > stop {
		status = "down"
	} else {
		status = "up"
	}

	priceAlert := domain.PriceAlert{
		Symbol:        coinSymbol,
		PriceAtMoment: coin.PriceUsd,
		StopLimit:     stopLimit,
		Status:        status,
	}

	s.UsersCache.SetPriceAlert(userId, priceAlert, time.Hour*48)
}

func (s *Service) RemovePriceAlert(userId string, coinSymbol string) {
	_, err := s.GetCoinInfo(coinSymbol)
	if err != nil {
		return
	}

	s.UsersCache.DeletePriceAlert(userId, coinSymbol)
}

func (s *Service) UpdateAlertsStatus() []domain.PriceAlert {

	var activeAlerts []domain.PriceAlert

	for id, user := range s.UsersCache.GetAllUsers() {
		for _, priceAlert := range user.PriceAlertsList {

			priceAtMoment, priceParseErr := strconv.ParseFloat(priceAlert.PriceAtMoment, 64)
			if priceParseErr != nil {
				return []domain.PriceAlert{}
			}

			stopLimit, stopLimitParseErr := strconv.ParseFloat(priceAlert.StopLimit, 64)
			if stopLimitParseErr != nil {
				return []domain.PriceAlert{}
			}

			if (priceAtMoment >= stopLimit && priceAlert.Status == "up") || (priceAtMoment <= stopLimit && priceAlert.Status == "down") {
				activeAlerts = append(activeAlerts, priceAlert)

				priceAlert.Status = "done"

				s.UsersCache.SetPriceAlert(id, priceAlert, time.Minute)
			}
		}
	}

	return activeAlerts
}

func (s *Service) AddCollectionItem(userId string, coinSymbol string) {
	coin, err := s.GetCoinInfo(coinSymbol)
	if err != nil {
		return
	}

	s.UsersCache.SetCollectionItem(userId, coin)
}

func (s *Service) RemoveCollectionItem(userId string, coinSymbol string) {
	coin, err := s.GetCoinInfo(coinSymbol)
	if err != nil {
		return
	}

	s.UsersCache.DeleteCollectionItem(userId, coin)
}

func (s *Service) GetCollection(userId string) []domain.Coin {
	return s.UsersCache.GetAllCollectionItems(userId)
}
