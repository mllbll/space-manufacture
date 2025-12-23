package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mllbll/space-manufacture/order/internal/model"
)

func (s *ServiceSuite) TestCancelError() {
	var (
		param = gofakeit.UUID()
		repoErr = gofakeit.Error()

		getOrderResponse = model.GetOrderResponce{
			OrderUUID: param,
			Status:    "PENDING_PAYMENT",
		}
	)
	
	s.orderRepository.On("Get", s.ctx, param).Return(getOrderResponse, nil)

	s.orderRepository.On("Cancel", s.ctx, param).Return(model.Order{},repoErr)

	err := s.service.Cancel(s.ctx, param)
	s.Require().Error(err)
	s.Require().Equal(err, repoErr)
}

func (s *ServiceSuite) TestCancelConflictError() {
	var (
		param = gofakeit.UUID()

		getOrderResponse = model.GetOrderResponce{
			OrderUUID: param,
			Status:    "PAID",
		}
	)
	
	s.orderRepository.On("Get", s.ctx, param).Return(getOrderResponse, nil)

//	s.orderRepository.On("Cancel", s.ctx, param).Return(model.Order{},model.ErrOrderConflict)

	err := s.service.Cancel(s.ctx, param)
	s.Require().Error(err)
	s.Require().Equal(err, model.ErrOrderConflict)
}

func (s *ServiceSuite) TestCancelNoFound() {
	var (
		param = gofakeit.UUID()

		getOrderResponse = model.GetOrderResponce{
			OrderUUID: param,
			Status:    "CANCELLED",
		}
	)
	
	s.orderRepository.On("Get", s.ctx, param).Return(getOrderResponse, nil)

//	s.orderRepository.On("Cancel", s.ctx, param).Return(model.Order{},model.ErrOrderNoContent)

	err := s.service.Cancel(s.ctx, param)
	s.Require().Error(err)
	s.Require().Equal(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelNoContent() {
	var (
		param = gofakeit.UUID()

		getOrderResponse = model.GetOrderResponce{
			OrderUUID: param,
			Status:    "PENDING_PAYMENT",
		}
	)
	
	s.orderRepository.On("Get", s.ctx, param).Return(getOrderResponse, nil)

	s.orderRepository.On("Cancel", s.ctx, param).Return(model.Order{},nil)

	err := s.service.Cancel(s.ctx, param)
	s.Require().Error(err)
	s.Require().Equal(err, model.ErrOrderNoContent)
}

func (s *ServiceSuite) TestCancelGetError() {
	var (
		param   = gofakeit.UUID()
		repoErr = gofakeit.Error()
	)

	s.orderRepository.On("Get", s.ctx, param).Return(model.GetOrderResponce{}, repoErr)

	err := s.service.Cancel(s.ctx, param)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
