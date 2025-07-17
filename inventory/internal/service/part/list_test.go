package part

import (
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *ServiceSuite) TestListEmptyFilter() {
	filter := model.PartsFilter{}
	s.partRepository.On("List", s.ctx, filter).Return([]model.Part{}, nil)
	parts, err := s.service.List(s.ctx, filter)

	require.NoError(s.T(), err)
	assert.Len(s.T(), parts, 3)
}

//func (s *RepositoryTestSuite) TestList_FilterByUUID() {
//	filter := model.PartsFilter{
//		Uuids: []string{"part-001", "part-003"},
//	}
//
//	parts, err := s.repo.List(s.ctx, filter)
//
//	require.NoError(s.T(), err)
//	assert.Len(s.T(), parts, 2)
//	assert.Equal(s.T(), "part-001", parts[0].Uuid)
//	assert.Equal(s.T(), "part-003", parts[1].Uuid)
//}
//
//func (s *RepositoryTestSuite) TestList_FilterByName() {
//	filter := model.PartsFilter{
//		Names: []string{"fuel tank", "WINDOW PORTHOLE"},
//	}
//
//	parts, err := s.repo.List(s.ctx, filter)
//
//	require.NoError(s.T(), err)
//	assert.Len(s.T(), parts, 2)
//	assert.Equal(s.T(), "Fuel Tank", parts[0].Name)
//	assert.Equal(s.T(), "Window Porthole", parts[1].Name)
//}
//
//func (s *RepositoryTestSuite) TestList_FilterByCategory() {
//	filter := model.PartsFilter{
//		Categories: []model.Category{model.ENGINE, model.FUEL},
//	}
//
//	parts, err := s.repo.List(s.ctx, filter)
//
//	require.NoError(s.T(), err)
//	assert.Len(s.T(), parts, 2)
//	assert.Equal(s.T(), model.ENGINE, parts[0].Category)
//	assert.Equal(s.T(), model.FUEL, parts[1].Category)
//}
//
//func (s *RepositoryTestSuite) TestList_FilterByCountry() {
//	filter := model.PartsFilter{
//		ManufacturerCountries: []string{"usa", "FRANCE"},
//	}
//
//	parts, err := s.repo.List(s.ctx, filter)
//
//	require.NoError(s.T(), err)
//	assert.Len(s.T(), parts, 2)
//	assert.Equal(s.T(), "USA", parts[0].Manufacturer.Country)
//	assert.Equal(s.T(), "France", parts[1].Manufacturer.Country)
//}
//
//func (s *RepositoryTestSuite) TestList_FilterByTags() {
//	filter := model.PartsFilter{
//		Tags: []string{"safety", "DURABLE"},
//	}
//
//	parts, err := s.repo.List(s.ctx, filter)
//
//	require.NoError(s.T(), err)
//	assert.Len(s.T(), parts, 2)
//	assert.Contains(s.T(), parts[0].Tags, "safety")
//	assert.Contains(s.T(), parts[1].Tags, "durable")
//}
//
//func (s *RepositoryTestSuite) TestList_CombinedFilters() {
//	filter := model.PartsFilter{
//		Categories:            []model.Category{model.ENGINE, model.FUEL},
//		ManufacturerCountries: []string{"USA"},
//		Tags:                  []string{"high-performance"},
//	}
//
//	parts, err := s.repo.List(s.ctx, filter)
//
//	require.NoError(s.T(), err)
//	assert.Len(s.T(), parts, 1)
//	assert.Equal(s.T(), "Engine Turbine", parts[0].Name)
//}
//
//func (s *RepositoryTestSuite) TestList_NoMatches() {
//	filter := model.PartsFilter{
//		Names: []string{"Non-existent Part"},
//	}
//
//	parts, err := s.repo.List(s.ctx, filter)
//
//	require.NoError(s.T(), err)
//	assert.Empty(s.T(), parts)
//}
