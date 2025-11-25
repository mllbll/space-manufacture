package service

import (
	"context"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, req model.GetPartRequest) (model.GetPartResponse, error)
	ListParts(ctx context.Context, req model.ListPartsRequest) (model.ListPartsResponse, error)
}


