package kafka

import "github.com/dfg007star/go_rocket/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
