package payment

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/payment/internal/model"
)

func (s *ServiceSuite) TestPayOrderSucccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID = gofakeit.UUID()
		paymentMethod = gofakeit.IntRange(0, 4)
		transactionUUID = gofakeit.UUID()

		payOrderMessage = model.PayOrderMessage{
			OrderUUID: orderUUID,
			UserUUID: userUUID,
			PaymentMethod: model.PaymentMethodEnum(paymentMethod),
		}

		payOrderRequest = model.PayOrderRequest{
			PayOrderMessage: payOrderMessage,
		}

		payOrderResponce = model.PayOrderResponse{
			TransactionUUID: transactionUUID,
		}
	)

	s.paymentRepository.On("PayOrder", s.ctx, payOrderRequest).Return(payOrderResponce, nil)

	res, err := s.service.PayOrder(s.ctx, payOrderRequest)
	s.NoError(err)
	s.Equal(payOrderResponce, res)
}

func (s *ServiceSuite) TestPayOrderError() {
	var (
		repoErr = gofakeit.Error()
		orderUUID = gofakeit.UUID()
		userUUID = gofakeit.UUID()
		paymentMethod = gofakeit.IntRange(0, 4)

		payOrderMessage = model.PayOrderMessage{
			OrderUUID: orderUUID,
			UserUUID: userUUID,
			PaymentMethod: model.PaymentMethodEnum(paymentMethod),
		}

		payOrderRequest = model.PayOrderRequest{
			PayOrderMessage: payOrderMessage,
		}
	)
	s.paymentRepository.On("PayOrder", s.ctx, payOrderRequest).Return(model.PayOrderResponse{}, repoErr)

	res, err := s.paymentRepository.PayOrder(s.ctx, payOrderRequest)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(res)
}
