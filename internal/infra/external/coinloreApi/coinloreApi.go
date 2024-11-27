package coinloreApi

import (
	"crypto-helper/internal/domain"
	"encoding/json"
	"net/http"
)

type Client struct{}

func (c *Client) GetCoinsInfo() (map[string]domain.Coin, error) {
	var data domain.Coins

	response, responseError := http.Get("https://api.coinlore.net/api/tickers/?start=0&limit=20")

	if responseError != nil {
		return map[string]domain.Coin{}, responseError
	}

	defer response.Body.Close()

	byteSlice := make([]byte, 102400)

	n, _ := response.Body.Read(byteSlice)

	unmarshallingError := json.Unmarshal(byteSlice[:n], &data)

	if unmarshallingError != nil {
		return map[string]domain.Coin{}, unmarshallingError
	}

	coins := data.Coins

	coinsMap := make(map[string]domain.Coin)

	for _, coin := range coins {
		coinsMap[coin.Symbol] = coin
	}

	return coinsMap, nil
}
