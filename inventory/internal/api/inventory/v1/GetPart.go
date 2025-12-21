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

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	parts, err := a.inventoryService.GetPart(ctx, converter.GetPartRequestToModel(req))

	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.Uuid)
		}
		return nil, err
	}
	return converter.GetPartResponceToProto(model.GetPartResponse{Parts: parts.Parts}), nil
}
