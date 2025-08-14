package kafka

import "github.com/dfg007star/go_rocket/notification/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
