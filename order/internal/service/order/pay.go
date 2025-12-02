package order

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/model"
)

func (s *service) Pay(ctx context.Context, param string, req model.PayOrderRequest) (model.PayOrderResponse, error) {
	order, err := s.orderRepository.Get(ctx, param)
	if err != nil {
		return model.PayOrderResponse{}, model.ErrOrderNotFound
	}
	// Идем в пеймент и получаем transaction_uuid
	transaction_uuid, err := s.paymentClient.PayOrder(ctx, order.OrderUUID, order.UserUUID, req)

	err = s.orderRepository.Pay(ctx, param, req, transaction_uuid)
	if err != nil {
		return model.PayOrderResponse{}, model.ErrOrderConflict
	}
	return model.PayOrderResponse{TransactionUUID: transaction_uuid}, nil
}
