package payment

import (
	"github.com/mllbll/space-manufacture/payment/internal/repository"
	def "github.com/mllbll/space-manufacture/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct {
	paymentRepository repository.PaymentRepository
}

func NewService(paymentRepository repository.PaymentRepository) *service {
	return &service{
		paymentRepository: paymentRepository,
	}
}
