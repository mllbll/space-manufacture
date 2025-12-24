package v1

import (
	"context"
	"testing"

	"github.com/mllbll/space-manufacture/order/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	orderService *mocks.OrderService

	api *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.orderService = mocks.NewOrderService(s.T())

	s.api = NewAPI(
		s.orderService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
