package domain

type User struct {
	UserId          string
	PriceAlertsList []PriceAlert
	Collection      []Coin
}

type PriceAlert struct {
	Symbol        string
	PriceAtMoment string
	StopLimit     string
	Status        string
}
