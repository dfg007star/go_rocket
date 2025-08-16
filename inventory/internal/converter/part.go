package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
)

func PartsModelToGrpcResponse(models []*model.Part) *inventoryV1.ListPartsResponse {
	if models == nil {
		return &inventoryV1.ListPartsResponse{
			Parts: nil,
		}
	}

	parts := make([]*inventoryV1.Part, 0, len(models))
	for _, m := range models {
		part := &inventoryV1.Part{
			Uuid:          m.Uuid,
			Name:          m.Name,
			Description:   m.Description,
			Price:         m.Price,
			StockQuantity: m.StockQuantity,
			Category:      inventoryV1.Category(m.Category),
			Dimensions: &inventoryV1.Dimensions{
				Length: m.Dimensions.Length,
				Width:  m.Dimensions.Width,
				Height: m.Dimensions.Height,
				Weight: m.Dimensions.Weight,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    m.Manufacturer.Name,
				Country: m.Manufacturer.Country,
				Website: m.Manufacturer.Website,
			},
			Tags:      m.Tags,
			Metadata:  convertMetadata(m.Metadata),
			CreatedAt: timestamppb.New(m.CreatedAt),
			UpdatedAt: timestamppb.New(m.UpdatedAt),
		}
		parts = append(parts, part)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: parts,
	}
}

func PartModelToGrpcResponse(model *model.Part) *inventoryV1.GetPartResponse {
	if model == nil {
		return nil
	}

	part := &inventoryV1.Part{
		Uuid:          model.Uuid,
		Name:          model.Name,
		Description:   model.Description,
		Price:         model.Price,
		StockQuantity: model.StockQuantity,
		Category:      inventoryV1.Category(model.Category),
		Dimensions: &inventoryV1.Dimensions{
			Length: model.Dimensions.Length,
			Width:  model.Dimensions.Width,
			Height: model.Dimensions.Height,
			Weight: model.Dimensions.Weight,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    model.Manufacturer.Name,
			Country: model.Manufacturer.Country,
			Website: model.Manufacturer.Website,
		},
		Tags:      model.Tags,
		Metadata:  convertMetadata(model.Metadata),
		CreatedAt: timestamppb.New(model.CreatedAt),
		UpdatedAt: timestamppb.New(model.UpdatedAt),
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}
}

func convertMetadata(metadata map[string]model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		result[key] = convertValue(value)
	}
	return result
}

func convertValue(value model.Value) *inventoryV1.Value {
	result := &inventoryV1.Value{}

	switch {
	case value.StringValue != nil:
		result.Value = &inventoryV1.Value_StringValue{
			StringValue: *value.StringValue,
		}
	case value.Int64Value != nil:
		result.Value = &inventoryV1.Value_Int64Value{
			Int64Value: *value.Int64Value,
		}
	case value.DoubleValue != nil:
		result.Value = &inventoryV1.Value_DoubleValue{
			DoubleValue: *value.DoubleValue,
		}
	case value.BoolValue != nil:
		result.Value = &inventoryV1.Value_BoolValue{
			BoolValue: *value.BoolValue,
		}
	}

	return result
}

func ConvertCategories(categories []inventoryV1.Category) []model.Category {
	if categories == nil {
		return nil
	}

	result := make([]model.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, model.Category(category))
	}

	return result
}
