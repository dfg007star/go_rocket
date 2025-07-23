package converter

import (
	"time"

	"github.com/dfg007star/go_rocket/order/internal/model"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
)

func PartsListToModel(parts []*inventoryV1.Part) []*model.Part {
	result := make([]*model.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, PartToModel(part))
	}

	return result
}

func PartToModel(part *inventoryV1.Part) *model.Part {
	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		tmp := part.UpdatedAt.AsTime()
		updatedAt = &tmp
	}

	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      (model.Category)(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      PartMetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     updatedAt,
	}
}

func DimensionsToModel(dimensions *inventoryV1.Dimensions) model.Dimensions {
	if dimensions == nil {
		return model.Dimensions{}
	}

	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) model.Manufacturer {
	if manufacturer == nil {
		return model.Manufacturer{}
	}

	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func PartMetadataToModel(metadata map[string]*inventoryV1.Value) map[string]model.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]model.Value, len(metadata))
	for key, value := range metadata {
		if value == nil {
			continue
		}

		result[key] = ValueToModel(value)
	}

	return result
}

func ValueToModel(value *inventoryV1.Value) model.Value {
	if value == nil {
		return model.Value{}
	}

	result := model.Value{}

	switch v := value.Value.(type) {
	case *inventoryV1.Value_StringValue:
		stringValue := v.StringValue
		result.StringValue = &stringValue
	case *inventoryV1.Value_Int64Value:
		int64Value := v.Int64Value
		result.Int64Value = &int64Value
	case *inventoryV1.Value_DoubleValue:
		doubleValue := v.DoubleValue
		result.DoubleValue = &doubleValue
	case *inventoryV1.Value_BoolValue:
		boolValue := v.BoolValue
		result.BoolValue = &boolValue
	}

	return result
}

func PartsFilterToProto(filter *model.PartsFilter) *inventoryV1.PartsFilter {
	if len(filter.Uuids) == 0 && len(filter.Names) == 0 && len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 && len(filter.Tags) == 0 {
		return nil
	}

	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
