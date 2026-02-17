package app

import (
	"context"

	"github.com/mllbll/space-manufacture/payment/internal/repository"
	"github.com/mllbll/space-manufacture/payment/internal/service"

	paymentV1API "github.com/mllbll/space-manufacture/payment/internal/api/payment/v1"
	paymentRepository "github.com/mllbll/space-manufacture/payment/internal/repository/payment"
	paymentService "github.com/mllbll/space-manufacture/payment/internal/service/payment"

	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1Api paymentV1.PaymentServiceServer

	paymentService service.PaymentService

	paymentRepository repository.PaymentRepository
}

func (d *diContainer) PaymentRepository(ctx context.Context) repository.PaymentRepository {
	if d.paymentRepository == nil {
		d.paymentRepository = paymentRepository.NewRepository()
	}

	return d.paymentRepository
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService(d.PaymentRepository(ctx))
	}

	return d.paymentService
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1Api == nil {
		d.paymentV1Api = paymentV1API.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1Api
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}
