package order

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	userUuid := gofakeit.UUID()
	orderUuid := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID()}

	var (
		name          = gofakeit.Name()
		description   = gofakeit.Paragraph(3, 5, 5, " ")
		price         = gofakeit.Price(100, 1000)
		stockQuantity = gofakeit.Int64()
		dimensions    = model.Dimensions{
			Height: gofakeit.Float64Range(1.0, 10.0),
			Width:  gofakeit.Float64Range(1.0, 10.0),
			Length: gofakeit.Float64Range(1.0, 10.0),
			Weight: gofakeit.Float64Range(0.1, 5.0),
		}
		manufacturer = model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}
		createdAt = time.Now()
	)

	tags := make([]string, gofakeit.Number(1, 5))
	for i := range tags {
		tags[i] = gofakeit.Word()
	}

	part := &model.Part{
		Uuid:          orderUuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      model.WING,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          tags,
		CreatedAt:     createdAt,
	}

	orderResp := &model.Order{
		OrderUuid:  orderUuid,
		UserUuid:   userUuid,
		PartUuids:  partUuids,
		TotalPrice: float32(price),
		Status:     model.PENDING_PAYMENT,
		CreatedAt:  createdAt,
	}

	filter := &model.PartsFilter{
		Uuids: partUuids,
	}

	listParts := []*model.Part{part}

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(listParts, nil).Once()
	s.orderRepository.On("Create", s.ctx, userUuid, listParts).Return(orderResp, nil).Once()
	resp, err := s.service.Create(s.ctx, &model.OrderCreate{UserUuid: userUuid, PartUuids: partUuids})

	s.NoError(err)
	s.Equal(orderResp, resp)
}

func (s *ServiceSuite) TestCreateOrderListPartsErr() {
	userUuid := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID()}

	filter := &model.PartsFilter{
		Uuids: partUuids,
	}

	expectedListPartsError := gofakeit.Error()

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(nil, expectedListPartsError).Once()
	resp, err := s.service.Create(s.ctx, &model.OrderCreate{UserUuid: userUuid, PartUuids: partUuids})

	s.Error(err)
	s.Empty(resp)
	s.Equal(err, expectedListPartsError)
}

func (s *ServiceSuite) TestCreateOrderErr() {
	userUuid := gofakeit.UUID()
	orderUuid := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID()}

	var (
		name          = gofakeit.Name()
		description   = gofakeit.Paragraph(3, 5, 5, " ")
		price         = gofakeit.Price(100, 1000)
		stockQuantity = gofakeit.Int64()
		category      = model.WING
		dimensions    = model.Dimensions{
			Height: gofakeit.Float64Range(1.0, 10.0),
			Width:  gofakeit.Float64Range(1.0, 10.0),
			Length: gofakeit.Float64Range(1.0, 10.0),
			Weight: gofakeit.Float64Range(0.1, 5.0),
		}
		manufacturer = model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}
		createdAt = time.Now()
	)

	part := &model.Part{
		Uuid:          orderUuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		CreatedAt:     createdAt,
	}

	part2 := &model.Part{
		Uuid:          orderUuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		CreatedAt:     createdAt,
	}

	filter := &model.PartsFilter{
		Uuids: partUuids,
	}

	listParts := []*model.Part{part, part2}
	expectedErr := model.ErrNotAllPartsMatched

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(listParts, nil).Once()
	resp, err := s.service.Create(s.ctx, &model.OrderCreate{UserUuid: userUuid, PartUuids: partUuids})

	s.Error(err)
	s.Empty(resp)
	s.Equal(err, expectedErr)
}

func (s *ServiceSuite) TestCreateOrderRepoErr() {
	userUuid := gofakeit.UUID()
	orderUuid := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID()}

	var (
		name          = gofakeit.Name()
		description   = gofakeit.Paragraph(3, 5, 5, " ")
		price         = gofakeit.Price(100, 1000)
		stockQuantity = gofakeit.Int64()
		category      = model.FUEL
		dimensions    = model.Dimensions{
			Height: gofakeit.Float64Range(1.0, 10.0),
			Width:  gofakeit.Float64Range(1.0, 10.0),
			Length: gofakeit.Float64Range(1.0, 10.0),
			Weight: gofakeit.Float64Range(0.1, 5.0),
		}
		manufacturer = model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}
		createdAt = time.Now()
	)

	part := &model.Part{
		Uuid:          orderUuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		CreatedAt:     createdAt,
	}

	filter := &model.PartsFilter{
		Uuids: partUuids,
	}

	listParts := []*model.Part{part}
	expectedErr := gofakeit.Error()

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(listParts, nil).Once()
	s.orderRepository.On("Create", s.ctx, userUuid, []*model.Part{part}).Return(&model.Order{}, expectedErr).Once()
	resp, err := s.service.Create(s.ctx, &model.OrderCreate{UserUuid: userUuid, PartUuids: partUuids})

	s.Error(err)
	s.Empty(resp)
	s.Equal(err, expectedErr)
}
