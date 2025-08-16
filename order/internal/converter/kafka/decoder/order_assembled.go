package decoder

import (
	"fmt"

	"github.com/dfg007star/go_rocket/order/internal/model"
	eventsV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUuid:    pb.EventUuid,
		OrderUuid:    pb.OrderUuid,
		UserUuid:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
