package order

import (
	"context"
	"log"
	"fmt"
	"github.com/mllbll/space-manufacture/order/internal/model"
	repoConverter "github.com/mllbll/space-manufacture/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, req model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

  if _, ok := r.data[req.OrderUUID]; ok {
	// уже существует
  log.Print("Заказ с таким UUID уже существует в системе")
  return fmt.Errorf("order with UUID %s already exists", req.OrderUUID)
	}

    // не существует — добавляем
  r.data[req.OrderUUID] = repoConverter.OrderToRepoModel(req)
	return nil
}
