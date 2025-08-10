package model

type ShipAssembledEvent struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

type OrderPaidEvent struct {
	EventUuid       string
	OrderUuid       string
	UserUuid        string
	TransactionUuid string
	PaymentMethod   string
}
