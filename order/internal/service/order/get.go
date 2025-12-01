package order

import (
	"context"

	"github.com/mllbll/space-manufacture/order/internal/model"
	repoConverter "github.com/mllbll/space-manufacture/order/internal/repository/converter"
)

func (s *service) Get(ctx context.Context, param string) (model.GetOrderResponce, error) {
	order, err := s.orderRepository.Get(ctx, param)
	if err != nil {
		return model.GetOrderResponce{}, err
	}

	return repoConverter.GetOrderResponseToModel(order), nil
}
