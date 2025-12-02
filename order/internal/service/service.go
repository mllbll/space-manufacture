package service

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/model"
)

type OrderService interface {
	Cancel(ctx context.Context, param string) error
	Get(ctx context.Context, params string) (model.GetOrderResponce, error)
	Pay(ctx context.Context, param string, req model.PayOrderRequest) (model.PayOrderResponse, error)
	Create(ctx context.Context, req model.CreateOrderRequest) (model.CreateOrderResponse, error)
}
