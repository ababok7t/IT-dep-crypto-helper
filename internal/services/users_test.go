package services

import (
	"crypto-helper/internal/domain"
	"crypto-helper/internal/infra/cache"
	"crypto-helper/internal/infra/external/coinloreApi"
	"math"
	"testing"
)

func TestAddCollectionItem(t *testing.T) {
	botService := Service{
		CoinsCache:   cache.NewCoinsCache(),
		UsersCache:   cache.NewUsersCache(),
		Client:       &coinloreApi.Client{},
		StateMachine: NewStateMachine(),
	}
	var list []domain.Coin
	expected :=
	circle := Circle{Radius: radius}
	actual, err := circle.GetArea()

	if err != nil {
		t.Error(err.Error())
	}

	if actual != expected {
		t.Error("incorrect result")
	}
}