package converter

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
)

func DimensionsToModel(info *inventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: info.Length,
		Width: info.Width,
		Height: info.Height,
		Weight: info.Weight,
	}
}

func DimensionsToProto(info model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: info.Length,
		Width: info.Width,
		Height: info.Height,
	}
}

func ManufacturerToModel(info *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name: info.Name,
		Country: info.Country,
		Website: info.Website,
	}
}

func ManufacturerToProto(info model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name: info.Name,
		Country: info.Country,
		Website: info.Website,
	}
}



// Конвертация слайса категорий из proto в model
func CategoriesSliceToModel(categories []inventoryV1.Category) []model.Category {
	result := make([]model.Category, len(categories))
	for i, cat := range categories {
		result[i] = model.Category(cat)
	}
	return result
}

// Конвертация слайса категорий из model в proto
func CategoriesSliceToProto(categories []model.Category) []inventoryV1.Category {
	result := make([]inventoryV1.Category, len(categories))
	for i, cat := range categories {
		result[i] = inventoryV1.Category(cat)
	}
	return result
}

func ValueToModel(val *inventoryV1.Value) model.Value {
	if val == nil {
		return nil
	}

	switch val.Types.(type) {
	case *inventoryV1.Value_StringValue:
		return model.StringValue{Value: val.GetStringValue()}
	case *inventoryV1.Value_Int64Value:
		return model.Int64Value{Value: val.GetInt64Value()}
	case *inventoryV1.Value_DoubleValue:
		return model.DoubleValue{Value: val.GetDoubleValue()}
	case *inventoryV1.Value_BoolValue:
		return model.BoolValue{Value: val.GetBoolValue()}
	default:
		return nil
	}
}

func ValueToProto(val model.Value) *inventoryV1.Value {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case model.StringValue :
		return &inventoryV1.Value{
			Types : &inventoryV1.Value_StringValue{
			StringValue: v.Value,
		},
	}
	case model.Int64Value :
		return &inventoryV1.Value{
			Types: &inventoryV1.Value_Int64Value{
				Int64Value: v.Value,
			},
		}
	case model.DoubleValue :
		return &inventoryV1.Value{
			Types: &inventoryV1.Value_DoubleValue{
				DoubleValue: v.Value,
			},
		}
	case model.BoolValue :
		return &inventoryV1.Value{
			Types: &inventoryV1.Value_BoolValue{
				BoolValue: v.Value,
			},
		}
	default :
		return nil
	}
}

func MetadataToModel(metadata map[string]*inventoryV1.Value) map[string]model.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]model.Value, len(metadata))
	for k, v := range metadata {
		result[k] = ValueToModel(v)
	}
	return result
}

func MetadataToProto(metadata map[string]model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*inventoryV1.Value, len(metadata))
	for k, v := range metadata {
		result[k] = ValueToProto(v)
	}
	return result
}

func PartsFilterToModel(info *inventoryV1.PartsFilter) model.PatrsFilter {
	return model.PatrsFilter{
		UUIDs: info.Uuids,
		Names: info.Names,
		Categories: CategoriesSliceToModel(info.Categories),
		Manufacture_contries: info.ManufacturerContries,
		Tags: info.Tags,
	}
}

func PartFilterToProto(info model.PatrsFilter) *inventoryV1.PartsFilter {
	return &inventoryV1.PartsFilter{
		Uuids: info.UUIDs,
		Names: info.Names,
		Categories: CategoriesSliceToProto(info.Categories),
		ManufacturerContries: info.Manufacture_contries,
		Tags: info.Tags,
	}
}

func PartToModel(info *inventoryV1.Part) model.Part {
	CreatedAt := lo.ToPtr(info.CreatedAt.AsTime())
	UpdatedAt := lo.ToPtr(info.UpdatedAt.AsTime())
	return model.Part{
		UUID: info.Uuid,
		Name: info.Name,
		Description: info.Description,
		Price: info.Price,
		Stock_quantity: info.StockQuantity,
		Category: model.Category(info.Category),
		Dimensions: DimensionsToModel(info.Dimensions),
		Manufacturer: ManufacturerToModel(info.Manufacturer),
		Tags: info.Tags,
		Metadata: MetadataToModel(info.Metadata),
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
	}
}

func PartToProto(info model.Part) *inventoryV1.Part {
	CreatedAt := timestamppb.New(*info.CreatedAt)
	UpdatedAt := timestamppb.New(*info.UpdatedAt)
	return &inventoryV1.Part {
		Uuid: info.UUID,
		Name: info.Name,
		Description: info.Description,
		Price: info.Price,
		StockQuantity: info.Stock_quantity,
		Category: inventoryV1.Category(info.Category),
		// нужно написать нормально Dimensions
		Dimensions: DimensionsToProto(info.Dimensions),
		Manufacturer: ManufacturerToProto(info.Manufacturer),
		Tags: info.Tags,
		Metadata: MetadataToProto(info.Metadata),
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
	}
}

func PartsSliceToModel(parts []*inventoryV1.Part) []model.Part {
	result := make([]model.Part, len(parts))
	for i, part := range parts {
		result[i] = PartToModel(part)
	}
	return result
}

func PartsSliceToProto(parts []model.Part) []*inventoryV1.Part {
	result := make([]*inventoryV1.Part, len(parts))
	for i, part := range parts {
		result[i] = PartToProto(part)
	}
	return result
}

func GetPartRequestToModel(info *inventoryV1.GetPartRequest) model.GetPartRequest {
	return model.GetPartRequest{
		UUID: info.Uuid,
	}
}

func GetPartRequestToProto(info model.GetPartRequest) *inventoryV1.GetPartRequest {
	return &inventoryV1.GetPartRequest{
		Uuid: info.UUID,
	}
}

func GetPartResponceToModel(info *inventoryV1.GetPartResponse) model.GetPartResponse {
	return model.GetPartResponse{
		Parts: PartToModel(info.Parts),
	}
}

func GetPartResponceToProto(info model.GetPartResponse) *inventoryV1.GetPartResponse {
	return &inventoryV1.GetPartResponse{
		Parts: PartToProto(info.Parts),
	}
}

func ListPartsRequestToModel(info *inventoryV1.ListPartsRequest) model.ListPartsRequest {
	return model.ListPartsRequest{
		Filter: PartsFilterToModel(info.Filter),
	}
}

func ListPartsRequestToProto(info model.ListPartsRequest) *inventoryV1.ListPartsRequest {
	return &inventoryV1.ListPartsRequest{
		Filter: PartFilterToProto(info.Filter),
	}
}

func ListPartsResponseToModel(info *inventoryV1.ListPartsResponse) model.ListPartsResponse {
	return model.ListPartsResponse{
		Parts: PartsSliceToModel(info.Parts),
	}
}

func ListPartsResponseToProto(info model.ListPartsResponse) *inventoryV1.ListPartsResponse {
	return &inventoryV1.ListPartsResponse{
		Parts: PartsSliceToProto(info.Parts),
	}
}

