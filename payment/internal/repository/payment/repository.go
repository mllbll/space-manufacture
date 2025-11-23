package payment

import (
	"sync"

	def "github.com/mllbll/space-manufacture/payment/internal/repository"
	repoModel "github.com/mllbll/space-manufacture/payment/internal/repository/model"
)

var _ def.PaymentRepository = (*repository)(nil)

type repository struct {
	mu sync.RWMutex
	data map[string]repoModel.PayOrderMessage
}

func NewRepository () *repository {
	return &repository{
		data: make(map[string]repoModel.PayOrderMessage),
	}
}


