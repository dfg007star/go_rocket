package part

import (
	"context"
	"strings"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, f *model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filter := converter.PartsFilterModelToPartsFilterRepoModel(f)

	parts := make([]*model.Part, 0, len(r.data))

	if isEmptyFilter(filter) {
		for _, part := range r.data {
			convertedPart := converter.RepoModelToPartModel(&part)
			parts = append(parts, convertedPart)
		}
		return parts, nil
	}

	// создаем map для более быстрого поиска значения :)
	uuidSet := makeStringSet(filter.Uuids)
	nameSet := makeStringSet(filter.Names)
	categorySet := makeCategorySet(filter.Categories)
	countrySet := makeStringSet(filter.ManufacturerCountries)
	tagSet := makeStringSet(filter.Tags)

	for _, part := range r.data {
		if isMatchAnyFilter(&part, uuidSet, nameSet, categorySet, countrySet, tagSet) {
			parts = append(parts, converter.RepoModelToPartModel(&part))
		}
	}

	return parts, nil
}

func isMatchAnyFilter(part *repoModel.Part,
	uuidSet map[string]struct{},
	nameSet map[string]struct{},
	categorySet map[repoModel.Category]struct{},
	countrySet map[string]struct{},
	tagSet map[string]struct{},
) bool {
	if len(uuidSet) > 0 {
		if _, exists := uuidSet[part.Uuid]; !exists {
			return false
		}
	}

	if len(nameSet) > 0 {
		if _, exists := nameSet[strings.ToLower(part.Name)]; !exists {
			return false
		}
	}

	if len(categorySet) > 0 {
		if _, exists := categorySet[part.Category]; !exists {
			return false
		}
	}

	if len(countrySet) > 0 {
		if _, exists := countrySet[strings.ToLower(part.Manufacturer.Country)]; !exists {
			return false
		}
	}

	if len(tagSet) > 0 {
		hasMatchingTag := false
		for _, tag := range part.Tags {
			if _, exists := tagSet[strings.ToLower(tag)]; exists {
				hasMatchingTag = true
				break
			}
		}
		if !hasMatchingTag {
			return false
		}
	}

	return true
}

func isEmptyFilter(filter *repoModel.PartsFilter) bool {
	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

func makeStringSet(items []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range items {
		set[strings.ToLower(item)] = struct{}{}
	}

	return set
}

func makeCategorySet(categories []repoModel.Category) map[repoModel.Category]struct{} {
	set := make(map[repoModel.Category]struct{})
	for _, cat := range categories {
		set[cat] = struct{}{}
	}

	return set
}
