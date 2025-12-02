package repository

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/model"
)

type OrderRepository interface {
	Cancel(ctx context.Context, param string) (model.Order, error)
	Get(ctx context.Context, params string) (model.GetOrderResponce, error)
	Pay(ctx context.Context, param string, req model.PayOrderRequest, transactionUUID string) error
	Create(ctx context.Context, req model.Order) error
}
