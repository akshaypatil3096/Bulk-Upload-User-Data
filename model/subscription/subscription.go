package model

type Subscription struct {
	Plan          string `json:"plan"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	Term          string `json:"term"`
}
