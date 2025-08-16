package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/dfg007star/go_rocket/notification/internal/model"
	eventsV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/events/v1"
)

type paidDecoder struct{}

func NewOrderPaidDecoder() *paidDecoder {
	return &paidDecoder{}
}

func (d *paidDecoder) PaidDecode(data []byte) (model.OrderPaidEvent, error) {
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
