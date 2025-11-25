package part

import (
	"context"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, req model.GetPartRequest) (model.GetPartResponse, error) {
	part, err := s.inventoryRepository.GetPart(ctx, req)

	if err != nil {
		return model.GetPartResponse{}, err
	}

	return part, nil
}
