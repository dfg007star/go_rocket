package model

import (
	"database/sql"
	"time"
)

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
	UpdatedAt       sql.NullTime
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
	}[s]
}
