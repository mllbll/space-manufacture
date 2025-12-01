package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mllbll/space-manufacture/order/internal/model"
)

func (s *service) Create(ctx context.Context, req model.CreateOrderRequest) (model.CreateOrderResponse, error) {
	newOrderUUID := uuid.NewString()

	getPriceMessage := model.PatrsFilter{
		UUIDs: req.PartUUIDs,
		Names: nil,
		Categories: nil,
		Manufacture_contries: nil,
		Tags: nil,
	}

	listParts, err := s.inventoryClient.ListParts(ctx, getPriceMessage)
	if err != nil {
		return model.CreateOrderResponse{}, model.ErrOrderNotFound
	}
	if len(listParts.Parts) != len(listParts.Parts) {
		return model.CreateOrderResponse{}, errors.New("Ошибка: не удалось получить всех деталей")
	}

	var sumPrice float64

	for _, part := range listParts.Parts {
		sumPrice += part.Price
	}

	newOrder := model.Order{
		OrderUUID: newOrderUUID,
		UserUUID: req.UserUUID,
		PartUUIDs: req.PartUUIDs,
		TotalPrice: float32(sumPrice),
		Status: "PENDING_PAYMENT",
	} 

	createSuccesStatus := s.orderRepository.Create(ctx, newOrder)
	if createSuccesStatus != nil {
		return model.CreateOrderResponse{}, createSuccesStatus
	}

	return model.CreateOrderResponse{
		OrderUUID: newOrderUUID,
		TotalPrice: float32(sumPrice),
	}, nil
}
