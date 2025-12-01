package v1

import (
	"context"

	clientConverter "github.com/mllbll/space-manufacture/order/internal/client/converter"
	"github.com/mllbll/space-manufacture/order/internal/model"
	generatedInventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PatrsFilter) (model.ListPartsResponse, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{Filter: clientConverter.PartFilterToProto(filter)})
	
	if err != nil {
		return model.ListPartsResponse{}, err
	}

	return clientConverter.ListPartsResponseToModel(parts), nil
}
