package converter

import (
	"github.com/mllbll/space-manufacture/payment/internal/model"
	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

func PayOrderMessageToModel(message *paymentV1.PayOrderMessage) model.PayOrderMessage {
	return model.PayOrderMessage{
		OrderUUID: message.OrderUuid,
		UserUUID: message.UserUuid,
		PaymentMethod: model.PaymentMethodEnum(message.PaymentMethod),
	}
}

func PayOrderMessageToProto(message model.PayOrderMessage) *paymentV1.PayOrderMessage {
	return &paymentV1.PayOrderMessage{
		OrderUuid: message.OrderUUID,
		UserUuid: message.UserUUID,
		PaymentMethod: paymentV1.PaymentMethodEnum(message.PaymentMethod),
	}
}

func PayOrderRequestToModel(req *paymentV1.PayOrderRequest) model.PayOrderRequest {
	return model.PayOrderRequest{
		PayOrderMessage: PayOrderMessageToModel(req.PayOrderMessage),
	}
}

func PayOrderRequestToProto(req model.PayOrderRequest) *paymentV1.PayOrderRequest {
	return &paymentV1.PayOrderRequest{
		PayOrderMessage: PayOrderMessageToProto(req.PayOrderMessage),
	}
}

func PayOrderResponceToModel(res *paymentV1.PayOrderResponse) model.PayOrderResponse {
	return model.PayOrderResponse{
		TransactionUUID: res.TransactionUuid,
	}
}

func PayOrderResponceToProto(res model.PayOrderResponse) *paymentV1.PayOrderResponse {
	return &paymentV1.PayOrderResponse{
		TransactionUuid: res.TransactionUUID,
	}
}
