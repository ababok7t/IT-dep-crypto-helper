package services

import (
	"crypto-helper/internal/infra/cache"
	"crypto-helper/internal/infra/external/coinloreApi"
)

type Service struct {
	CoinsCache   *cache.CoinsCache
	UsersCache   *cache.UsersCache
	Client       *coinloreApi.Client
	StateMachine *StateMachine
}
