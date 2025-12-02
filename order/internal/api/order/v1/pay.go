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

func (a *api) PayOrder(ctx context.Context, req *generatedOrder.PayOrderRequest, params generatedOrder.PayOrderParams,) (generatedOrder.PayOrderRes, error) {
	order, err := a.orderService.Pay(ctx, params.OrderUUID, converter.PayOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, "Order no found")
		}
		return nil, err
	}

	return converter.PayOrderResponseToOpenAPI(order), nil
}


