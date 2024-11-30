package services

import (
	"crypto-helper/internal/domain"
	"testing"
)

func TestMakeCoinReply(t *testing.T) {
	coin := domain.Coin{
		Id:               "90",
		Symbol:           "BTC",
		Name:             "Bitcoin",
		PriceUsd:         "85000",
		PercentChange24H: "1",
		PercentChange1H:  "2",
		PercentChange7D:  "3",
	}
	actual := domain.MakeCoinReply("BTC", coin, "100000", "%")

	expected := "информация о криптовалюте BTC:\nназвание: Bitcoin\nцена: 85000$\nизменение цены за 1 час: 2 %\nизменение цены за 24 часа: 1 %\nизменение цены за 7 дней: 3 %\nпрогноз: 100000 $"

	if actual != expected {
		t.Error("incorrect result")
	}
}
