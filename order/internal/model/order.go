package model

import "time"

type PaymentMethod int

const (
	UNSPECIFIED PaymentMethod = iota
	CARD
	SBP
	CREDIT_CARD
	INVESTOR_MONEY
)

type Status int

const (
	PENDING_PAYMENT Status = iota
	PAID
	CANCELLED
	COMPLETED
)

type Order struct {
	OrderUuid       string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float32
	TransactionUuid *string
	PaymentMethod   *PaymentMethod
	Status          Status
	CreatedAt       time.Time
}

type OrderCreate struct {
	UserUuid  string
	PartUuids []string
}

type OrderCreateInfo struct {
	OrderUuid  string
	TotalPrice float32
}

type OrderUpdate struct {
	OrderUuid       string
	TransactionUuid *string
	PaymentMethod   *PaymentMethod
	Status          *Status
}

func (pm PaymentMethod) String() string {
	return [...]string{
		"UNSPECIFIED",
		"CARD",
		"SBP",
		"CREDIT_CARD",
		"INVESTOR_MONEY",
	}[pm]
}

func (s Status) String() string {
	return [...]string{
		"PENDING_PAYMENT",
		"PAID",
		"CANCELLED",
		"COMPLETED",
	}[s]
}
