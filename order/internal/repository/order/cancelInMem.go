package order

import (
	"context"

	repoConverter "github.com/mllbll/space-manufacture/order/internal/repository/converter"
	"github.com/mllbll/space-manufacture/order/internal/model"
)

func (r *repository) Cancel(ctx context.Context, param string) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[param]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	order.Status = "CANCELLED"
	r.data[param] = order

	return repoConverter.OrderToModel(order), nil
}
