package kafka

import "github.com/dfg007star/go_rocket/notification/internal/model"

type OrderPaidDecoder interface {
	PaidDecode(data []byte) (model.OrderPaidEvent, error)
}

type OrderAssembledDecoder interface {
	AssembledDecode(data []byte) (model.ShipAssembledEvent, error)
}
