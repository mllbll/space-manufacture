package order

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/model"
	"github.com/samber/lo"
)

func (r *repository) Pay(ctx context.Context, param string, req model.PayOrderRequest, transactionUUID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[param]
	if !ok {
		return model.ErrOrderNotFound
	}

	order.Status = "PAID"
	order.TransactionUUID = lo.ToPtr(transactionUUID)
	order.PaymentMethod = lo.ToPtr(int32(req.PaymentMethod))

	r.data[param] = order

	return nil

}
