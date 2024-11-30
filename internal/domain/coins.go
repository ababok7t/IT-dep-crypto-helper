package domain

import "fmt"

type Coin struct {
	Id               string `json:"id"`
	Symbol           string `json:"symbol"`
	Name             string `json:"name"`
	PriceUsd         string `json:"price_usd"`
	PercentChange24H string `json:"percent_change_24h"`
	PercentChange1H  string `json:"percent_change_1h"`
	PercentChange7D  string `json:"percent_change_7d"`
}

type Coins struct {
	Coins []Coin `json:"data"`
}

func MakeCoinReply(coinSymbol string, coinInfo Coin, coinForecast string, pr string) string {
	reply := fmt.Sprintf("информация о криптовалюте %s:\nназвание: %s\nцена: %s$\nизменение цены за 1 час: %s %s\nизменение цены за 24 часа: %s %s\nизменение цены за 7 дней: %s %s\nпрогноз: %s $", coinSymbol, coinInfo.Name, coinInfo.PriceUsd, coinInfo.PercentChange1H, pr, coinInfo.PercentChange24H, pr, coinInfo.PercentChange7D, pr, coinForecast)
	return reply
}
