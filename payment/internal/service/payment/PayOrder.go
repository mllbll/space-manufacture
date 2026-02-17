package payment

import (
	"context"
	"log"
	"platform/pkg/logger"

	"github.com/mllbll/space-manufacture/payment/internal/model"
	"go.uber.org/zap"
)

func (s *service) PayOrder(ctx context.Context, message model.PayOrderRequest) (model.PayOrderResponse, error) {
	res, err := s.paymentRepository.PayOrder(ctx, message)
	if err != nil {
		logger.Error(ctx, "Ошибка Оплаты в репо слое", zap.Error(err))
		return model.PayOrderResponse{}, err
	}

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", res.TransactionUUID)

	return res, nil
}
