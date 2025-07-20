package model

type PaymentMethod int

const (
	PAYMENT_METHOD_UNSPECIFIED PaymentMethod = iota
	PAYMENT_METHOD_CARD
	PAYMENT_METHOD_SBP
	PAYMENT_METHOD_CREDIT_CARD
	PAYMENT_METHOD_INVESTOR_MONEY
)

func (p PaymentMethod) String() string {
	switch p {
	case PAYMENT_METHOD_CARD:
		return "PAYMENT_METHOD_CARD"
	case PAYMENT_METHOD_SBP:
		return "PAYMENT_METHOD_SBP"
	case PAYMENT_METHOD_CREDIT_CARD:
		return "PAYMENT_METHOD_CREDIT_CARD"
	case PAYMENT_METHOD_INVESTOR_MONEY:
		return "PAYMENT_METHOD_INVESTOR_MONEY"
	default:
		return "PAYMENT_METHOD_UNSPECIFIED"
	}
}

type Payment struct {
	OrderUuid       string
	UserUuid        string
	TransactionUuid string
	PaymentMethod   PaymentMethod
}
