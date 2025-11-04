package main

import (
	"context"
	"sync"
	"time"

	orderV1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
	"golang.org/x/mod/sumdb/storage"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.GetOrderResponse
}

//мапка ордер хранит UUID и структуру GetOrderResponse
//но нужно нормально обернуть или сделать в опенапи декларации новую структуру Order потому что хранить в гетордер это не вайб

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.GetOrderResponse),
	}
}

func (s *OrderStorage) GetOrder(order_uuid string) *orderV1.GetOrderResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) CreateOrder(order *orderV1.GetOrderResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order_uuid := "3f3b7af8-5646-4d11-b373-10c01d6f9c05" //тут тоже заглушка ибо я хз откуда я должен вообще высрать этот UUID
	s.orders[order_uuid] = order
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) GetOrderByUuid(_ context.Context, params orderV1.APIV1OrdersOrderUUIDGetParams) (orderV1.APIV1OrdersOrderUUIDGetRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by UUID " + params.OrderUUID + " not found",
		}, nil
	}

	return order, nil
}

// Заглушил поля кроме нужных в реквесте приколами, нужно исправить!!!
func (h *OrderHandler) CreateNewOrder(_ context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderResponse, error) {
	order := &orderV1.GetOrderResponse{
		OrderUUID:       "05b4fe37-7822-4d95-8f1f-76edbcc3c134",
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      12,
		TransactionUUID: "6cfb5601-43d1-4ad0-8d26-b0c96840a7e3",
		PaymentMethod:   orderV1.NewOptGetOrderResponsePaymentMethod(1), //Я вообще хз что тут за тип данных, скорее всего наебнулась генерация
		Status:          "payed",
	}
	order_resp := &orderV1.CreateOrderResponse{
		OrderUUID:  orderV1.NewOptString("05b4fe37-7822-4d95-8f1f-76edbcc3c134"),
		TotalPrice: orderV1.NewOptFloat32(12.1),
	}
	h.storage.CreateOrder(order)

	return *order_resp, nil
}
