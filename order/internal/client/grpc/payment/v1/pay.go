package v1

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/client/converter"
	"github.com/mllbll/space-manufacture/order/internal/model"
	generatedPaymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID string, req model.PayOrderRequest) (string, error) {
	payOrderMessage := &generatedPaymentV1.PayOrderMessage{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		// Костыль с конвертором
		PaymentMethod: generatedPaymentV1.PaymentMethodEnum(generatedPaymentV1.PaymentMethodEnum_value[converter.PaymentMethodEnumToString(req.PaymentMethod)]),
	}
	
	res, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{PayOrderMessage: payOrderMessage})

	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}


