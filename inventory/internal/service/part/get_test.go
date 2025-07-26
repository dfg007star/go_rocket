package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/stretchr/testify/require"
)

func (s *ServiceSuite) TestGet() {
	var (
		partUUID     = gofakeit.UUID()
		expectedPart = model.Part{
			Uuid:          partUUID,
			Name:          "Engine Turbine",
			Category:      model.ENGINE,
			Price:         9999.99,
			StockQuantity: 5,
			Dimensions: model.Dimensions{
				Length: 2.5,
				Weight: 150.3,
			},
			Manufacturer: model.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
			},
			CreatedAt: time.Now(),
		}
	)

	s.partRepository.On("Get", s.ctx, partUUID).Return(expectedPart, nil)

	result, err := s.service.Get(s.ctx, partUUID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedPart, result)
	s.partRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetNotFound() {
	partUUID := gofakeit.UUID()
	expectedError := model.ErrPartNotFound

	s.partRepository.On("Get", s.ctx, partUUID).Return(model.Part{}, expectedError)

	result, err := s.service.Get(s.ctx, partUUID)
	require.Error(s.T(), err)
	require.Empty(s.T(), result.Uuid)
	require.ErrorIs(s.T(), err, expectedError)
	s.partRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetEmptyUUID() {
	_, err := s.service.Get(s.ctx, "")
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "uuid cannot be empty")
}
