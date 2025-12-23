package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/order/internal/model"
	"github.com/samber/lo"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		param = gofakeit.UUID()

		orderUUID = gofakeit.UUID()
		userUUID = gofakeit.UUID()
		partUUIDs = []string{gofakeit.UUID()}
		totalPrice = gofakeit.Float32()
		transactionUUID = lo.ToPtr(gofakeit.UUID())
		paymentMethod = lo.ToPtr(model.PAYMENT_METHOD_ENUM_CARD)
		status = gofakeit.Word()

		getOrderResponse = model.GetOrderResponce{
			OrderUUID: orderUUID,
			UserUUID: userUUID,
			PartUUIDs: partUUIDs,
			TotalPrice: totalPrice,
			TransactionUUID: transactionUUID,
			PaymentMethod: paymentMethod,
			Status: status,
		}
	)

	s.orderRepository.On("Get", s.ctx, param).Return(getOrderResponse, nil)

	res, err := s.service.Get(s.ctx, param)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(res, getOrderResponse)
}

func (s *ServiceSuite) TestGetError() {
	var (
		param = gofakeit.UUID()
		repoErr = gofakeit.Error()
	)

	s.orderRepository.On("Get", s.ctx, param).Return(model.GetOrderResponce{}, repoErr)

	res, err := s.service.Get(s.ctx, param)

	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().Equal(err, repoErr)
}


func (s *ServiceSuite) TestGetNotFoundError() {
	var (
		param = gofakeit.UUID()
	)

	s.orderRepository.On("Get", s.ctx, param).Return(model.GetOrderResponce{}, model.ErrOrderNotFound)

	res, err := s.service.Get(s.ctx, param)

	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().Equal(err, model.ErrOrderNotFound)
}
