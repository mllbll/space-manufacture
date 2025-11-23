package converter

import (
	"github.com/mllbll/space-manufacture/payment/internal/model"
	repoModel "github.com/mllbll/space-manufacture/payment/internal/repository/model"
)

func PayOrderMessageToRepoModel(message model.PayOrderMessage) repoModel.PayOrderMessage {
	return repoModel.PayOrderMessage{
		OrderUUID:     message.OrderUUID,
		UserUUID:      message.UserUUID,
		PaymentMethod: repoModel.PaymentMethodEnum(message.PaymentMethod), // конвертация типов
	}
}

func PayOrderMessageToModel(message repoModel.PayOrderMessage) model.PayOrderMessage {
	return model.PayOrderMessage{
		OrderUUID: message.OrderUUID,
		UserUUID: message.UserUUID,
		PaymentMethod: model.PaymentMethodEnum(message.PaymentMethod),
	} // конвертация типов
}

func PayOrderRequestToRepoModel(req model.PayOrderRequest) repoModel.PayOrderRequest {
	return repoModel.PayOrderRequest{
		PayOrderMessage: PayOrderMessageToRepoModel(req.PayOrderMessage),
	}
}

func PayOrderRequestToModel(req repoModel.PayOrderRequest) model.PayOrderRequest {
	return model.PayOrderRequest{
		PayOrderMessage: PayOrderMessageToModel(req.PayOrderMessage),
	}
}

func PayOrderResponseToRepoModel(res model.PayOrderResponse) repoModel.PayOrderResponse {
	return repoModel.PayOrderResponse{
		TransactionUUID: res.TransactionUUID,
	}
}

func PayOrderResponseToModel (res repoModel.PayOrderResponse) model.PayOrderResponse {
	return model.PayOrderResponse{
		TransactionUUID: res.TransactionUUID,
	}
}
