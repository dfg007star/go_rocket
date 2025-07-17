package part

import (
	"context"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/converter"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/mocks"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	partRepository *mocks.PartRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.partRepository = mocks.NewPartRepository(s.T())

	s.service = NewService(
		s.partRepository,
	)

	data := []repoModel.Part{
		{
			Uuid:     "part-001",
			Name:     "Engine Turbine",
			Category: repoModel.ENGINE,
			Manufacturer: repoModel.Manufacturer{
				Country: "USA",
				Name:    "Amazon",
			},
			Dimensions: repoModel.Dimensions{
				Weight: 8.0,
			},
			Tags: []string{"critical", "high-performance"},
		},
		{
			Uuid:     "part-002",
			Name:     "Fuel Tank",
			Category: repoModel.FUEL,
			Manufacturer: repoModel.Manufacturer{
				Country: "Germany",
				Name:    "Amazon",
			},
			Dimensions: repoModel.Dimensions{
				Weight: 12.0,
			},
			Tags: []string{"storage", "safety"},
		},
		{
			Uuid:     "part-003",
			Name:     "Window Porthole",
			Category: repoModel.PORTHOLE,
			Manufacturer: repoModel.Manufacturer{
				Country: "France",
				Name:    "Amazon",
			},
			Dimensions: repoModel.Dimensions{
				Weight: 9.0,
			},
			Tags: []string{"view", "durable"},
		},
	}

	for _, part := range data {
		_, err := s.service.Create(s.ctx, converter.RepoModelToPartModel(&part))
		if err != nil {
			continue
		}
	}
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
