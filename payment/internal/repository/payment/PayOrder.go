package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/mllbll/space-manufacture/payment/internal/model"
	repoModel "github.com/mllbll/space-manufacture/payment/internal/repository/model"
	repoConverter "github.com/mllbll/space-manufacture/payment/internal/repository/converter"
)

func (r *repository) PayOrder(_ context.Context, message model.PayOrderRequest) (model.PayOrderResponse, error) {
	newTransactionUUID := uuid.NewString()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[newTransactionUUID] = repoModel.PayOrderMessage{
		OrderUUID: message.PayOrderMessage.OrderUUID,
		UserUUID: message.PayOrderMessage.UserUUID,
		PaymentMethod: repoModel.PaymentMethodEnum(message.PayOrderMessage.PaymentMethod),
	}

	return repoConverter.PayOrderResponseToModel(repoModel.PayOrderResponse{TransactionUUID: newTransactionUUID}), nil
}


