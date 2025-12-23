package order

import (
	"context"
	"testing"

	"github.com/mllbll/space-manufacture/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"

	grpcMock "github.com/mllbll/space-manufacture/order/internal/client/grpc/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *mocks.OrderRepository

	paymentClient *grpcMock.PaymentClient

	inventoryClient *grpcMock.InventoryClient

	service *service
}

func (s *ServiceSuite) SetupTest () {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.paymentClient = grpcMock.NewPaymentClient(s.T())
	s.inventoryClient = grpcMock.NewInventoryClient(s.T())

	s.service = NewService(
		s.orderRepository,
		s.paymentClient,
		s.inventoryClient,
		)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration (t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

