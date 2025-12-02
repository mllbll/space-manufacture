package part

import (
	"sync"

	def "github.com/mllbll/space-manufacture/inventory/internal/repository"
	repoModel "github.com/mllbll/space-manufacture/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data : make(map[string]repoModel.Part),
	}
}
