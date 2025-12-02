package service

import (
	"context"

	"github.com/mllbll/space-manufacture/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, message model.PayOrderRequest) (model.PayOrderResponse, error)
}
