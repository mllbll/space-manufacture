package order

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, param string) error {
	order, err := s.orderRepository.Cancel(ctx, param)
	if err != nil {
		return  err
	}

	switch order.Status {
	case "PENDING_PAYMENT":
		order.Status = "PAID"
		return model.ErrOrderNoContent
	case "PAID":
		return model.ErrOrderConflict
	}
	return model.ErrOrderNotFound
}


