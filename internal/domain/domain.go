package domain

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

type User struct {
	UserId          string
	PriceAlertsList []PriceAlert
}

type PriceAlert struct {
	Symbol        string
	PriceAtMoment string
	StopLimit     string
	status        string
}
