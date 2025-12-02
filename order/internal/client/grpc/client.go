package grpc

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PatrsFilter) (model.ListPartsResponse, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, req model.PayOrderRequest) (string, error)
}
