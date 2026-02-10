package part

import (
	"context"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, req model.ListPartsRequest) (model.ListPartsResponse, error) {
	return s.inventoryRepository.ListParts(ctx, req)
}
