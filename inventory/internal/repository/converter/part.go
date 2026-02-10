package converter

import (
	"github.com/mllbll/space-manufacture/inventory/internal/model"
	repoModel "github.com/mllbll/space-manufacture/inventory/internal/repository/model"
)

func DimensionsToRepoModel(info model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: info.Length,
		Width: info.Width,
		Height: info.Height,
		Weight: info.Weight,
	}
}

func DimensionsToModel(info repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: info.Length,
		Width: info.Width,
		Height: info.Height,
		Weight: info.Weight,
	}
}

func ManufacturerToRepoModel(info model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name: info.Name,
		Country: info.Country,
		Website: info.Website,
	}
}

func ManufacturerToModel(info repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name: info.Name,
		Country: info.Country,
		Website: info.Website,
	}
}

// Инициализируем новый слайс и пробрасываем в него значения
func CategoriesSliceToRepoModel(categories []model.Category) []repoModel.Category {
	result := make([]repoModel.Category, len(categories))
	for i, cat := range categories {
		result[i] = repoModel.Category(cat)
	}
	return result
}

// Инициализируем новый слайс и пробрасываем в него значения
func CategoriesSliceToModel(categories []repoModel.Category) []model.Category {
	result := make([]model.Category, len(categories))
	for i, cat := range categories {
		result[i] = model.Category(cat)
	}
	return result
}

func PartsFilterToRepoModel(info model.PatrsFilter) repoModel.PatrsFilter {
    return repoModel.PatrsFilter{
        UUIDs: info.UUIDs,
        Names: info.Names,
        Categories: CategoriesSliceToRepoModel(info.Categories),
        Manufacture_contries: info.Manufacture_contries,
        Tags: info.Tags,
    }
}

func PartsFilterToModel(info repoModel.PatrsFilter) model.PatrsFilter {
    return model.PatrsFilter{
        UUIDs: info.UUIDs,
        Names: info.Names,
        Categories: CategoriesSliceToModel(info.Categories),
        Manufacture_contries: info.Manufacture_contries,
        Tags: info.Tags,
    }
}

func ValueToRepoModel(val model.Value) repoModel.Value {
	switch v := val.(type) {
	case model.StringValue:
		return repoModel.StringValue{Value: v.Value}
	case model.Int64Value:
		return repoModel.Int64Value{Value: v.Value}
	case model.DoubleValue:
		return repoModel.DoubleValue{Value: v.Value}
	case model.BoolValue:
		return repoModel.BoolValue{Value: v.Value}
	default:
		return nil
	}
}

func ValueToModel(val repoModel.Value) model.Value {
	switch v := val.(type) {
	case repoModel.StringValue:
		return model.StringValue{Value: v.Value}
	case repoModel.Int64Value:
		return model.Int64Value{Value: v.Value}
	case repoModel.DoubleValue:
		return model.DoubleValue{Value: v.Value}
	case repoModel.BoolValue:
		return model.BoolValue{Value: v.Value}
	default:
		return nil
	}
}

// func MetadataToRepoModel(metadata map[string]model.Value) map[string]repoModel.Value {
// 	if metadata == nil {
// 		return nil
// 	}
// 	result := make(map[string]repoModel.Value, len(metadata))
// 	for k, v := range metadata {
// 		result[k] = ValueToRepoModel(v)
// 	}
// 	return result
// }
//
// func MetadataToModel(metadata map[string]repoModel.Value) map[string]model.Value {
// 	if metadata == nil {
// 		return nil
// 	}
// 	result := make(map[string]model.Value, len(metadata))
// 	for k, v := range metadata {
// 		result[k] = ValueToModel(v)
// 	}
// 	return result
// }

func MetadataToModel (m map[string]interface{}) map[string]model.Value {
	if m == nil {
		return nil 
	}
	result := make(map[string]model.Value, len(m))

	for k ,v := range m {
		if v == nil {
			continue
		}

		switch val := v.(type) {
		case string:
			result[k] = model.StringValue{Value: val}
		case int:
			result[k] = model.Int64Value{Value: int64(val)}
		case int32:
			result[k] = model.Int64Value{Value: int64(val)}
		case int64:
			result[k] = model.Int64Value{Value: val}
		case bool:
			result[k] = model.BoolValue{Value: val}
		case float64:
			if val == float64(int64(val)) {
				result[k] = model.Int64Value{Value: int64(val)}
			} else {
				result[k] = model.DoubleValue{Value: val}
			}
		case float32:
			result[k] = model.DoubleValue{Value: float64(val)}
		}
	}
	return result
}

func MetadataToRepoModel (m map[string]model.Value) map[string]interface{} {
	if m == nil {
		return nil
	}

	result := make(map[string]interface{}, len(m))
	
	for k, v := range m {
		if v == nil {
			continue
		}

		switch val := v.(type) {
		case model.Int64Value:
			result[k] = val.Value
		case model.BoolValue:
			result[k] = val.Value
		case model.DoubleValue:
			result[k] = val.Value
		case model.StringValue:
			result[k] = val.Value
		}
	}
	return result
}

func PartToRepoModel(info model.Part) repoModel.Part {
	return repoModel.Part{
		UUID: info.UUID,
		Name: info.Name,
		Description: info.Description,
		Price: info.Price,
		Stock_quantity: info.Stock_quantity,
		Category: int32(info.Category),
		Dimensions: DimensionsToRepoModel(info.Dimensions),
		Manufacturer: ManufacturerToRepoModel(info.Manufacturer),
		Tags: info.Tags,
		Metadata: MetadataToRepoModel(info.Metadata),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}

func PartToModel(info repoModel.Part) model.Part {
	return model.Part{
		UUID: info.UUID,
		Name: info.Name,
		Description: info.Description,
		Price: info.Price,
		Stock_quantity: info.Stock_quantity,
		Category: model.Category(info.Category),
		Dimensions: DimensionsToModel(info.Dimensions),
		Manufacturer: ManufacturerToModel(info.Manufacturer),
		Tags: info.Tags,
		Metadata: MetadataToModel(info.Metadata),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}

// Инициализируем новый слайс и пробрасываем в него значения перебирая их по индексу в поданном слайсе
func PartsSliceToRepoModel(parts []model.Part) []repoModel.Part {
	result := make([]repoModel.Part, len(parts))
	for i, cat := range parts {
		result[i] = PartToRepoModel(cat)
	}
	return result
}

// Инициализируем новый слайс и пробрасываем в него значения
func PartsSliceToModel(parts []repoModel.Part) []model.Part {
	result := make([]model.Part, len(parts))
	for i, cat := range parts {
		result[i] = PartToModel(cat)
	}
	return result
}

func GetPartRequestToRepoModel(info model.GetPartRequest) repoModel.GetPartRequest {
	return repoModel.GetPartRequest{
		UUID: info.UUID,
	}
}

func GetPartRequestToModel(info repoModel.GetPartRequest) model.GetPartRequest {
	return model.GetPartRequest{
		UUID: info.UUID,
	}
}

func GetPartResponseToRepoModel(info model.GetPartResponse) repoModel.GetPartResponse {
	return repoModel.GetPartResponse{
		Parts: PartToRepoModel(info.Parts),
	}
}

func GetPartResponseToModel(info repoModel.GetPartResponse) model.GetPartResponse {
	return model.GetPartResponse{
		Parts: PartToModel(info.Parts),
	}
}

func ListPartsRequestToRepoModel(info model.ListPartsRequest) repoModel.ListPartsRequest {
	return repoModel.ListPartsRequest{
		Filter: PartsFilterToRepoModel(info.Filter),
	}
}

func ListPartsRequestToModel(info repoModel.ListPartsRequest) model.ListPartsRequest {
	return model.ListPartsRequest{
		Filter: PartsFilterToModel(info.Filter),
	}
}

func ListPartsResponseToRepoModel(info model.ListPartsResponse) repoModel.ListPartsResponse {
	return repoModel.ListPartsResponse{
		Parts: PartsSliceToRepoModel(info.Parts),
	}
}

func ListPartsResponseToModel(info repoModel.ListPartsResponse) model.ListPartsResponse {
	return model.ListPartsResponse{
		Parts: PartsSliceToModel(info.Parts),
	}
}
