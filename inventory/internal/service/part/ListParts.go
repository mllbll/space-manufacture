package part

import (
	"context"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
)

// ищет элемент в срезе строк
func containsString(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// ищет элемент типа Category в срезе таких значений
func containsCategory(list []model.Category, cat model.Category) bool {
	for _, v := range list {
		if v == cat {
			return true
		}
	}
	return false
}

// проверяет есть ли хотя бы одно совпадение в двух списках
func hasOverlap(a, b []string) bool {
	m := make(map[string]struct{}, len(a))
	for _, v := range a {
		m[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := m[v]; ok {
			return true
		}
	}
	return false
}

func (s *service) ListParts(ctx context.Context, req model.ListPartsRequest) (model.ListPartsResponse, error) {
	resp, err := s.inventoryRepository.ListParts(ctx, req)
	if err != nil {
		return model.ListPartsResponse{}, err
	}

	if !isEmptyFilter(req.Filter) {
		resp.Parts = filterParts(resp.Parts, req.Filter)
	}

	return resp, nil

}

func isEmptyFilter(filter model.PatrsFilter) bool {
	return len(filter.UUIDs) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.Manufacture_contries) == 0 &&
		len(filter.Tags) == 0
}

func filterParts(parts []model.Part, filter model.PatrsFilter) []model.Part {
	var result []model.Part

	for _, part := range parts {
		if len(filter.UUIDs) > 0 && !containsString(filter.UUIDs, part.UUID) {
			continue
		}

		if len(filter.Names) > 0 && !containsString(filter.Names, part.Name) {
			continue
		}

		if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
			continue
		}

		if len(filter.Manufacture_contries) > 0 {
			if part.Manufacturer.Country == "" || !containsString(filter.Manufacture_contries, part.Manufacturer.Country) {
				continue
			}
		}

		if len(filter.Tags) > 0 && !hasOverlap(filter.Tags, part.Tags) {
			continue
		}
		result = append(result, part)
	}
	return result
}
