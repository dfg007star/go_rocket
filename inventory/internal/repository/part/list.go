package part

import (
	"context"
	"fmt"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *repository) List(ctx context.Context, f *model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filter := bson.M{}

	if len(f.Uuids) > 0 {
		filter["uuid"] = bson.M{"$in": f.Uuids}
	}

	if len(f.Names) > 0 {
		nameRegex := make([]bson.M, len(f.Names))
		for i, name := range f.Names {
			nameRegex[i] = bson.M{"name": primitive.Regex{Pattern: name, Options: "i"}}
		}
		filter["$or"] = nameRegex
	}

	if len(f.Categories) > 0 {
		filter["category"] = bson.M{"$in": f.Categories}
	}

	if len(f.ManufacturerCountries) > 0 {
		filter["manufacturer.country"] = bson.M{"$in": f.ManufacturerCountries}
	}

	if len(f.Tags) > 0 {
		filter["tags"] = bson.M{"$in": f.Tags}
	}

	cursor, err := r.data.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query parts: %w", err)
	}
	defer func() {
		err = cursor.Close(ctx)
		if err != nil {
			panic(err)
		}
	}()

	var parts []*repoModel.Part
	if err := cursor.All(ctx, &parts); err != nil {
		return nil, fmt.Errorf("failed to decode parts: %w", err)
	}

	result := make([]*model.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, converter.RepoModelToPartModel(part))
	}

	return result, nil
}
