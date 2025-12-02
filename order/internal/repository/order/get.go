package order

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/model"
	repoConvertet "github.com/mllbll/space-manufacture/order/internal/repository/converter"
	repoModel "github.com/mllbll/space-manufacture/order/internal/repository/model"
)

func (r *repository) Get(_ context.Context, params string) (model.GetOrderResponce, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[params]
	if !ok {
		return model.GetOrderResponce{}, model.ErrOrderNotFound
	}

	resp := repoModel.GetOrderResponce{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		Status: order.Status,
		PaymentMethod: order.PaymentMethod,
	}

	return repoConvertet.GetOrderResponseToModel(resp), nil
}
