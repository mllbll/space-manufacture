package v1

import (
	"context"

	"github.com/mllbll/space-manufacture/payment/internal/converter"
	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

func (a* api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	res, err := a.paymentService.PayOrder(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		return nil, err
	}


	return converter.PayOrderResponceToProto(res), nil
}


