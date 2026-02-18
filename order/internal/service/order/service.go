package order

import (
	"github.com/mllbll/space-manufacture/order/internal/client/grpc"
	"github.com/mllbll/space-manufacture/order/internal/repository"
	def "github.com/mllbll/space-manufacture/order/internal/service"
)

var _ def.OrderService = (*service)(nil)


type service struct {
	orderRepository repository.OrderRepository

	paymentClient grpc.PaymentClient

	inventoryClient grpc.InventoryClient
}

func NewService(
	orderRepository repository.OrderRepository,
	paymentClient grpc.PaymentClient,
	inventoryClient grpc.InventoryClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		paymentClient: paymentClient,
		inventoryClient: inventoryClient,
	}
}
