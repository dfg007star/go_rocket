package converter

import (
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
)

func modelCategoryToRepoCategory(category model.Category) repoModel.Category {
	switch category {
	case model.ENGINE:
		return repoModel.ENGINE
	case model.FUEL:
		return repoModel.FUEL
	case model.PORTHOLE:
		return repoModel.PORTHOLE
	case model.WING:
		return repoModel.WING
	default:
		return repoModel.UNKNOWN
	}
}

func repoCategoryToModelCategory(category repoModel.Category) model.Category {
	switch category {
	case repoModel.ENGINE:
		return model.ENGINE
	case repoModel.FUEL:
		return model.FUEL
	case repoModel.PORTHOLE:
		return model.PORTHOLE
	case repoModel.WING:
		return model.WING
	default:
		return model.UNKNOWN
	}
}

func modelDimensionsToRepoDimensions(dimensions model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func repoDimensionsToModelDimensions(dimensions repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func repoManufacturerToModelManufacturer(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func modelManufacturerToRepoManufacturer(manufacturer model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func repoMetadataToModelMetadata(metadata map[string]repoModel.Value) map[string]model.Value {
	result := make(map[string]model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = model.Value{
			StringValue: value.StringValue,
			Int64Value:  value.Int64Value,
			DoubleValue: value.DoubleValue,
			BoolValue:   value.BoolValue,
		}
	}
	return result
}

func modelMetadataToRepoMetadata(metadata map[string]model.Value) map[string]repoModel.Value {
	result := make(map[string]repoModel.Value, len(metadata))
	for key, value := range metadata {
		result[key] = repoModel.Value{
			StringValue: value.StringValue,
			Int64Value:  value.Int64Value,
			DoubleValue: value.DoubleValue,
			BoolValue:   value.BoolValue,
		}
	}
	return result
}

func PartToRepoModel(part *model.Part) *repoModel.Part {
	return &repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      modelCategoryToRepoCategory(part.Category),
		Dimensions:    modelDimensionsToRepoDimensions(part.Dimensions),
		Manufacturer:  modelManufacturerToRepoManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      modelMetadataToRepoMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func RepoModelToPartModel(repo *repoModel.Part) *model.Part {
	return &model.Part{
		Uuid:          repo.Uuid,
		Name:          repo.Name,
		Description:   repo.Description,
		Price:         repo.Price,
		StockQuantity: repo.StockQuantity,
		Category:      repoCategoryToModelCategory(repo.Category),
		Dimensions:    repoDimensionsToModelDimensions(repo.Dimensions),
		Manufacturer:  repoManufacturerToModelManufacturer(repo.Manufacturer),
		Tags:          repo.Tags,
		Metadata:      repoMetadataToModelMetadata(repo.Metadata),
		CreatedAt:     repo.CreatedAt,
		UpdatedAt:     repo.UpdatedAt,
	}
}
