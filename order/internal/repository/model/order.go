package model

type PaymentMethod int

const (
	UNSPECIFIED PaymentMethod = iota
	CARD
	SBP
	CREDIT_CARD
	INVESTOR_MONEY
)

type Order struct {
	OrderUuid       string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float32
	TransactionUuid *string
	PaymentMethod   *PaymentMethod
	Status          string
}
