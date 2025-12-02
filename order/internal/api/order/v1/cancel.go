package v1

import (
	"context"

	generatedOrder "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params generatedOrder.CancelOrderParams) (generatedOrder.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	return nil, err
}


