package payment

import (
	"context"
	"log"

	"github.com/mllbll/space-manufacture/payment/internal/model"
)

func(s *service) PayOrder(ctx context.Context, message model.PayOrderRequest) (model.PayOrderResponse, error) {
	res, err := s.PayOrder(ctx, message)
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", res.TransactionUUID)

	return res, nil
}
