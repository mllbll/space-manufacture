package v1

import (
	"context"
	"errors"
	"github.com/mllbll/space-manufacture/inventory/internal/converter"
	"github.com/mllbll/space-manufacture/inventory/internal/model"
	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.inventoryService.ListParts(ctx, converter.ListPartsRequestToModel(req))

	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with your filter not found")
		}
		return nil, err
	}

	return converter.ListPartsResponseToProto(model.ListPartsResponse{Parts: parts.Parts}), nil
}
