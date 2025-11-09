package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	orderV1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
	order_v1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

//мапка ордер хранит UUID и структуру GetOrderResponse
//но нужно нормально обернуть или сделать в опенапи декларации новую структуру Order потому что хранить в гетордер это не вайб
//Вроде обернул, но почему то проблемы с видением структуры ордер

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrderStorage) GetOrder(order_uuid string) *orderV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		return nil
	}

	return order
}

func newNotFound(message string) *orderV1.NotFoundError {
	return &orderV1.NotFoundError{
		Code:    404,
		Message: message,
	}
}

func (s *OrderStorage) CreateOrder(order *orderV1.CreateOrderRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// тут нет обработки на поиск order по UUID потому что мы его создаем
	//	order_uuid := "3f3b7af8-5646-4d11-b373-10c01d6f9c05" //тут тоже заглушка ибо я хз откуда я должен вообще высрать этот UUID
	// по идее UUID я должен парсить из созданной структуры
	order_uuid := "62e69f5b-9c60-4017-b095-97ce32e27042" // UUID все еще заглушил
	new_order := &orderV1.Order{
		OrderUUID:  order_uuid,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: 0.0, // тут тоже заглушка в виде 0 потому что от сервиса другого должна приходить общая стоимость
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	s.orders[order_uuid] = new_order
}

func (s *OrderStorage) PayOrder(order_uuid string, payment_method *orderV1.PayOrderRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//	paymentMethod := orderV1.PayOrderRequestPaymentMethod(payment_method)

	//	s.orders[order_uuid].PaymentMethod = orderV1.NewOptOrderPaymentMethod(payment_method.PaymentMethod)
	order, ok := s.orders[order_uuid]
	if !ok {
		return
	}
	// парсим поле PaymentMethod из PayOrderRequest и запихиваем его в структуру Order
	method, ok := payment_method.GetPaymentMethod().Get()
	if !ok {
		return // тут нужно написать обработчик ошибочек
	}
	// проверяем что order[order_uuid] существует и закидываем его в order
	order.PaymentMethod = orderV1.NewOptOrderPaymentMethod(orderV1.OrderPaymentMethod(method))
}

func (s *OrderStorage) CancelOrder(order_cancel_uuid *orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order_uuid := order_cancel_uuid.OrderUUID

	order, ok := s.orders[order_uuid]
	if !ok {
		return newNotFound(fmt.Sprintf("Order with UUID '%s' not found", order_uuid)), nil
	}

	order_status := order.Status
	if order_status == orderV1.OrderStatusPENDINGPAYMENT {
		order.Status = orderV1.OrderStatusCANCELLED
	}
	if order_status == orderV1.OrderStatusPAID {
		return &orderV1.CancelOrderConflict{}, nil
	}
	return nil, nil
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
	//	order := &orderV1.Order{ // Тут хз что должно быть потому что приходит реквест в котором мало полей
	//		OrderUUID:       "05b4fe37-7822-4d95-8f1f-76edbcc3c134",
	//		UserUUID:        req.UserUUID,
	//		PartUuids:       req.PartUuids,
	//		TotalPrice:      12,
	//		TransactionUUID: "6cfb5601-43d1-4ad0-8d26-b0c96840a7e3",
	// NewOptOrderPaymentMethod это функция которая принимает инт и преобразует его в нужный нам енум
	//		PaymentMethod: orderV1.NewOptOrderPaymentMethod(1), //Я вообще хз что тут за тип данных, скорее всего наебнулась генерация
	//		Status:        "PENDING_PAYMENT",
	//	}
	order := &orderV1.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUuids: req.PartUuids,
	}
	order_resp := &orderV1.CreateOrderResponse{
		OrderUUID:  orderV1.NewOptString("05b4fe37-7822-4d95-8f1f-76edbcc3c134"), // Заглушил значение UUID до момента пока не напишу норм функцию генерации UUID
		TotalPrice: orderV1.NewOptFloat32(12.1),
	}
	h.storage.CreateOrder(order)

	return *order_resp, nil
}

// адски насрал в PayOrder и в storage тут
func (h *OrderHandler) PayOrderByUUID(_ context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	orderPayInfo := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.NewOptPayOrderRequestPaymentMethod(1),
	}

	h.storage.PayOrder(params.OrderUUID, orderPayInfo)
	// прокидываем в storage наш s из функции PayOrder

	order_pay_resp := &order_v1.PayOrderResponse{
		TransactionUUID: "7b5c38cb-57b9-4f2e-9a0a-c518add9ccaa",
		// тут должен наверное генерироваться uuid транзакции
	}
	return order_pay_resp, nil
}

func (h *OrderHandler) CancelOrderByUUID(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order_uuid := &orderV1.CancelOrderParams{
		OrderUUID: params.OrderUUID,
	}

	h.storage.CancelOrder(order_uuid)

	return &orderV1.CancelOrderNoContent{}, nil
}
