package v1

import (
	"context"
	"errors"

	"github.com/mllbll/space-manufacture/order/internal/converter"
	"github.com/mllbll/space-manufacture/order/internal/model"
	generatedOrder "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Get(ctx context.Context, req *generatedOrder.APIV1OrdersOrderUUIDGetParams) (generatedOrder.APIV1OrdersOrderUUIDGetRes, error) {
	order, err := a.orderService.Get(ctx, req.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, "Order no found")
		}
		return nil, err
	}

	orderOpenAPI := converter.GetOrderResponseToOpenAPI(order)
	return orderOpenAPI, nil
}


