package part

import (
	"context"
	"platform/pkg/logger"

	"github.com/mllbll/space-manufacture/inventory/internal/model"
	"go.uber.org/zap"
)

func (s *service) GetPart(ctx context.Context, req model.GetPartRequest) (model.GetPartResponse, error) {
	part, err := s.inventoryRepository.GetPart(ctx, req)

	if err != nil {
		logger.Error(ctx, "Ошибка при получении", zap.Error(err))
		return model.GetPartResponse{}, err
	}

	return part, nil
}
