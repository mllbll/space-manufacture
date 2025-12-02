package repository

import (
	"context"
	"github.com/mllbll/space-manufacture/payment/internal/model"
)

type PaymentRepository interface {
	PayOrder(ctx context.Context, message model.PayOrderRequest) (model.PayOrderResponse, error)
}
