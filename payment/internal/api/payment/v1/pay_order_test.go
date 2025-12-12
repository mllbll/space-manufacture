package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"

	"github.com/mllbll/space-manufacture/payment/internal/converter"
	"github.com/mllbll/space-manufacture/payment/internal/model"
)

func (s *APISuite) TestPayOrderSuccess() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		paymentMethod   = gofakeit.IntRange(0, 4)
		transactionUUID = gofakeit.UUID()

		payOrderMessage = &paymentV1.PayOrderMessage{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethodEnum(paymentMethod),
		}

		payOrderRequest = &paymentV1.PayOrderRequest{
			PayOrderMessage: payOrderMessage,
		}

		payOrderResponce = &paymentV1.PayOrderResponse{
			TransactionUuid: transactionUUID,
		}

		expectedPayOrderRequest = converter.PayOrderRequestToModel(payOrderRequest)

		expectedPayOrderResponce = converter.PayOrderResponceToModel(payOrderResponce)
	)

	s.paymentService.On("PayOrder", s.ctx, expectedPayOrderRequest).Return(expectedPayOrderResponce, nil)

	res, err := s.api.PayOrder(s.ctx, payOrderRequest)
	s.Require().NoError(err)
	s.Require().Equal(payOrderResponce, res)
	s.Require().NotNil(res)
}

func (s *APISuite) TestPayOrderErr() {
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = gofakeit.IntRange(0, 4)
		expectedErr   = gofakeit.Error()

		payOrderMessage = &paymentV1.PayOrderMessage{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethodEnum(paymentMethod),
		}

		payOrderRequest = &paymentV1.PayOrderRequest{
			PayOrderMessage: payOrderMessage,
		}

		expectedPayOrderRequest = converter.PayOrderRequestToModel(payOrderRequest)
	)

	s.paymentService.On("PayOrder", s.ctx, expectedPayOrderRequest).Return(model.PayOrderResponse{}, expectedErr)

	res, err := s.api.PayOrder(s.ctx, payOrderRequest)

	s.Require().Error(err)
	s.Require().Empty(res)

	s.Require().ErrorIs(err, expectedErr)

}
