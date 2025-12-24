package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/order/internal/converter"
	"github.com/mllbll/space-manufacture/order/internal/model"
	generatedOrder "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestGetSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID = gofakeit.UUID()
		partUUIDs = []string{gofakeit.UUID()}
		totalPrice = gofakeit.Float32()
		transactionUUID = generatedOrder.OptString{}
		status = generatedOrder.OrderStatus("PAID")
		paymentMethod = generatedOrder.NewOptOrderPaymentMethod(generatedOrder.OrderPaymentMethod(gofakeit.Number(0,4)))

		params = &generatedOrder.APIV1OrdersOrderUUIDGetParams{
			OrderUUID: orderUUID,
		}

		getOrderRes = &generatedOrder.Order{
			OrderUUID: orderUUID,
			UserUUID: userUUID,
			PartUuids: partUUIDs,
			TotalPrice: totalPrice,
			TransactionUUID: transactionUUID,
			Status: status,
			PaymentMethod: paymentMethod,
		}

		expectedGetOrderRes = converter.GetOrderResponseToModel(getOrderRes)
	)

	s.orderService.On("Get", s.ctx, orderUUID).Return(expectedGetOrderRes, nil)

	res, err := s.api.APIV1OrdersOrderUUIDGet(s.ctx, *params)

	s.Require().NoError(err)
	s.Require().NotEmpty(res)
	s.Require().Equal(res, getOrderRes)
}

func (s *APISuite) TestGetError() {
	var (
		serviceErr = gofakeit.Error()
		orderUUID = gofakeit.UUID()

		params = &generatedOrder.APIV1OrdersOrderUUIDGetParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("Get", s.ctx, orderUUID).Return(model.GetOrderResponce{}, serviceErr)

	res, err := s.api.APIV1OrdersOrderUUIDGet(s.ctx, *params)

	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().Equal(err, serviceErr)
}
