package decoder

import (
	"fmt"

	"github.com/dfg007star/go_rocket/assembly/internal/model"
	eventsV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		EventUuid:       pb.EventUuid,
		OrderUuid:       pb.OrderUuid,
		UserUuid:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod,
		TransactionUuid: pb.TransactionUuid,
	}, nil
}
