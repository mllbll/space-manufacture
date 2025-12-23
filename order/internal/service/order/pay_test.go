package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/order/internal/model"
	"github.com/samber/lo"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		param = gofakeit.UUID()

		paymentMethodPayOrderRequest = model.PAYMENT_METHOD_ENUM_CARD

		transactionUUIDNoPointer = gofakeit.UUID()

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


		payOrderRequest = model.PayOrderRequest{
			PaymentMethod: paymentMethodPayOrderRequest,
		}

		payOrderResponse = model.PayOrderResponse{
			TransactionUUID: transactionUUIDNoPointer,
		}
	)

	s.orderRepository.On("Get", s.ctx, param).Return(getOrderResponse, nil)

	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, payOrderRequest).Return(transactionUUIDNoPointer, nil)

	s.orderRepository.On("Pay", s.ctx, param, payOrderRequest, transactionUUIDNoPointer).Return(nil)

	res, err := s.service.Pay(s.ctx, param, payOrderRequest)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(res, payOrderResponse)
}

func (s *ServiceSuite) TestPayError() {
	var (
		param = gofakeit.UUID()

		paymentMethodPayOrderRequest = model.PAYMENT_METHOD_ENUM_CARD


		payOrderRequest = model.PayOrderRequest{
			PaymentMethod: paymentMethodPayOrderRequest,
		}
	)
	s.orderRepository.On("Get", s.ctx, param).Return(model.GetOrderResponce{}, model.ErrOrderNotFound)

	res, err := s.service.Pay(s.ctx, param, payOrderRequest)

	s.Require().Error(err)
	s.Empty(res)
	s.Equal(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestPayClientError() {
	var (
		param = gofakeit.UUID()

		clientErr = gofakeit.Error()

		paymentMethodPayOrderRequest = model.PAYMENT_METHOD_ENUM_CARD

	//	transactionUUIDNoPointer = gofakeit.UUID()


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


		payOrderRequest = model.PayOrderRequest{
			PaymentMethod: paymentMethodPayOrderRequest,
		}
	)
	s.orderRepository.On("Get", s.ctx, param).Return(getOrderResponse, nil)

	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, payOrderRequest).Return("", clientErr)

	//s.orderRepository.On("Pay", s.ctx, param, payOrderRequest, transactionUUIDNoPointer).Return(nil)

	res, err := s.service.Pay(s.ctx, param, payOrderRequest)

	s.Require().Error(err)
	s.Empty(res)
	s.Equal(err, clientErr)
}

func (s *ServiceSuite) TestPayPayError() {
	var (
		param = gofakeit.UUID()

		paymentMethodPayOrderRequest = model.PAYMENT_METHOD_ENUM_CARD

		transactionUUIDNoPointer = gofakeit.UUID()

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

		payOrderRequest = model.PayOrderRequest{
			PaymentMethod: paymentMethodPayOrderRequest,
		}
	)
	s.orderRepository.On("Get", s.ctx, param).Return(getOrderResponse, nil)

	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, payOrderRequest).Return(transactionUUIDNoPointer, nil)

	s.orderRepository.On("Pay", s.ctx, param, payOrderRequest, transactionUUIDNoPointer).Return(model.ErrOrderConflict)

	res, err := s.service.Pay(s.ctx, param, payOrderRequest)

	s.Require().Error(err)
	s.Empty(res)
	s.Equal(err, model.ErrOrderConflict)
}


