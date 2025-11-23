package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/google/uuid"

	orderV1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"

	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"

	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	// Порты gRPC сервисов
	paymentServiceServerAddress   = "localhost:50051"
	inventoryServiceServerAddress = "localhost:50052"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

// Функция которая будет ходить в payment и получать TransactionUUID
func payOrderCall(req *paymentV1.PayOrderRequest) (string, error) {

	if req.PayOrderMessage == nil {
		return "", fmt.Errorf("pay_order_message is required")
	}

	ctx := context.Background()

	conn, err := grpc.NewClient(
		paymentServiceServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return "", fmt.Errorf("failed to connect to payment service: %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	client := paymentV1.NewPaymentServiceClient(conn)

	payOrderMessage := &paymentV1.PayOrderMessage{
		OrderUuid:     req.PayOrderMessage.OrderUuid,
		UserUuid:      req.PayOrderMessage.UserUuid,
		PaymentMethod: req.PayOrderMessage.PaymentMethod,
	}

	resp, err := client.PayOrder(ctx, &paymentV1.PayOrderRequest{PayOrderMessage: payOrderMessage})
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}

// convertPaymentMethodToOrder преобразует строковый PayOrderRequestPaymentMethod в integer OrderPaymentMethod
func convertPaymentMethodToOrder(method orderV1.PayOrderRequestPaymentMethod) orderV1.OrderPaymentMethod {
	switch method {
	case orderV1.PayOrderRequestPaymentMethodUNKNOWN:
		return orderV1.OrderPaymentMethod0
	case orderV1.PayOrderRequestPaymentMethodCARD:
		return orderV1.OrderPaymentMethod1
	case orderV1.PayOrderRequestPaymentMethodSBP:
		return orderV1.OrderPaymentMethod2
	case orderV1.PayOrderRequestPaymentMethodCREDITCARD:
		return orderV1.OrderPaymentMethod3
	case orderV1.PayOrderRequestPaymentMethodINVESTORMONEY:
		return orderV1.OrderPaymentMethod4
	default:
		return orderV1.OrderPaymentMethod0
	}
}

// convertPaymentMethodToEnum преобразует строковый PayOrderRequestPaymentMethod в PaymentMethodEnum для payment сервиса
func convertPaymentMethodToEnum(method orderV1.PayOrderRequestPaymentMethod) paymentV1.PaymentMethodEnum {
	switch method {
	case orderV1.PayOrderRequestPaymentMethodUNKNOWN:
		return paymentV1.PaymentMethodEnum_PAYMENT_METHOD_ENUM_UNSPECIFIED
	case orderV1.PayOrderRequestPaymentMethodCARD:
		return paymentV1.PaymentMethodEnum_PAYMENT_METHOD_ENUM_CARD
	case orderV1.PayOrderRequestPaymentMethodSBP:
		return paymentV1.PaymentMethodEnum_PAYMENT_METHOD_ENUM_SBP
	case orderV1.PayOrderRequestPaymentMethodCREDITCARD:
		return paymentV1.PaymentMethodEnum_PAYMENT_METHOD_ENUM_CREDIT_CARD
	case orderV1.PayOrderRequestPaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethodEnum_PAYMENT_METHOD_ENUM_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethodEnum_PAYMENT_METHOD_ENUM_UNSPECIFIED
	}
}

func inventoryCall(listOfParts []string) (float64, error) {
	//	if req.Filter == nil {
	//		return 0.0, fmt.Errorf("Filter is required")
	//	}

	ctx := context.Background()

	conn, err := grpc.NewClient(
		inventoryServiceServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return 0.0, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	client := inventoryV1.NewInventoryServiceClient(conn)

	// мб стоит заменить значения на нули тут
	getPriceMessage := &inventoryV1.PartsFilter{
		Uuids:                listOfParts,
		Names:                nil,
		Categories:           nil,
		ManufacturerContries: nil,
		Tags:                 nil,
	}

	listParts, err := client.ListParts(ctx, &inventoryV1.ListPartsRequest{Filter: getPriceMessage})
	if err != nil {
		log.Printf("Не удалось получить детали")
	}
	if len(listParts.Parts) != len(listOfParts) {
		log.Printf("Ошибка: не удалось получить всех деталей")
	}

	var sumPrice float64

	for _, part := range listParts.Parts {
		sumPrice += part.Price
	}

	return sumPrice, nil
}

//мапка ордер хранит UUID и структуру GetOrderResponse
//но нужно нормально обернуть или сделать в опенапи декларации новую структуру Order потому что хранить в гетордер это не вайб
//Вроде обернул, но почему то проблемы с видением структуры ордер

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrderStorage) GetOrder(order_uuid string) (*orderV1.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		//		return nil, newNotFound(fmt.Sprint("Order with UUID '%s' not found", order_uuid)), nil
		return nil, fmt.Errorf("order with UUID %q not found", order_uuid)
	}

	return order, nil
}

func newNotFound(message string) *orderV1.NotFoundError {
	return &orderV1.NotFoundError{
		Code:    404,
		Message: message,
	}
}

func (s *OrderStorage) CreateOrderByUUID(order_uuid string, order *orderV1.CreateOrderRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// тут нет обработки на поиск order по UUID потому что мы его создаем
	//	order_uuid := "3f3b7af8-5646-4d11-b373-10c01d6f9c05" //тут тоже заглушка ибо я хз откуда я должен вообще высрать этот UUID
	// по идее UUID я должен парсить из созданной структуры
	//	order_uuid := "62e69f5b-9c60-4017-b095-97ce32e27042" // UUID все еще заглушил
	//	order_uuid := uuid.New().String()

	totalPrice, err := inventoryCall(order.PartUuids)
	if err != nil {
		log.Printf("Ошибка получения общей суммы")
	}

	//	order_uuid := uuid.New().String()

	new_order := &orderV1.Order{
		OrderUUID:  order_uuid,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: float32(totalPrice), // идем в inventoryService и получаем список деталей
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	s.orders[order_uuid] = new_order

	// return user_uuid
}

func (s *OrderStorage) PayOrderByUUID(order_uuid string, payment_method *orderV1.PayOrderRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	//	paymentMethod := orderV1.PayOrderRequestPaymentMethod(payment_method)

	//	s.orders[order_uuid].PaymentMethod = orderV1.NewOptOrderPaymentMethod(payment_method.PaymentMethod)
	order, ok := s.orders[order_uuid]
	if !ok {
		//		return newNotFound(fmt.Sprint("Order with UUID 's' not found", order_uuid))
		return fmt.Errorf("order with UUID %q not found", order_uuid)
	}
	// парсим поле PaymentMethod из PayOrderRequest и запихиваем его в структуру Order
	method, ok := payment_method.GetPaymentMethod().Get()
	if !ok {
		return errors.New("payment method is required")
		// тут нужно написать обработчик ошибочек
	}
	// Преобразуем строковый PayOrderRequestPaymentMethod в integer OrderPaymentMethod
	orderPaymentMethod := convertPaymentMethodToOrder(method)
	// проверяем что order[order_uuid] существует и закидываем его в order
	order.PaymentMethod = orderV1.NewOptOrderPaymentMethod(orderPaymentMethod)
	order.Status = orderV1.OrderStatusPAID

	return nil
}

func (s *OrderStorage) CancelOrderByUUID(order_cancel_uuid *orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
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
	return &orderV1.CancelOrderNoContent{}, nil
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) APIV1OrdersOrderUUIDGet(_ context.Context, params orderV1.APIV1OrdersOrderUUIDGetParams) (orderV1.APIV1OrdersOrderUUIDGetRes, error) {
	order, _ := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order by UUID " + params.OrderUUID + " not found",
		}, nil
	}

	return order, nil
}

// Заглушил поля кроме нужных в реквесте приколами, нужно исправить!!!
func (h *OrderHandler) AddNewOrder(_ context.Context, req *orderV1.CreateOrderRequest) (orderV1.AddNewOrderRes, error) {
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

	totalPrice, err := inventoryCall(req.PartUuids)
	if err != nil {
		log.Printf("Ошибка получение общей стоимости")
	}

	order_uuid := uuid.New().String()
	order_resp := &orderV1.CreateOrderResponse{
		//		OrderUUID:  orderV1.NewOptString("05b4fe37-7822-4d95-8f1f-76edbcc3c134"), // Заглушил значение UUID до момента пока не напишу норм функцию генерации UUID
		OrderUUID:  orderV1.NewOptString(order_uuid),
		TotalPrice: orderV1.NewOptFloat32(float32(totalPrice)),
	}
	h.storage.CreateOrderByUUID(order_uuid, order)

	return order_resp, nil
}

// адски насрал в PayOrder и в storage тут
func (h *OrderHandler) PayOrder(_ context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	// Проверяем наличие payment method
	method, ok := req.GetPaymentMethod().Get()
	if !ok {
		return nil, fmt.Errorf("payment method is required")
	}

	orderPayInfo := &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.NewOptPayOrderRequestPaymentMethod(method),
	}

	//	h.storage.PayOrderByUUID(params.OrderUUID, orderPayInfo
	if err := h.storage.PayOrderByUUID(params.OrderUUID, orderPayInfo); err != nil {
		return nil, err
	}
	// прокидываем в storage наш s из функции PayOrder

	// Получаем заказ для извлечения UserUUID
	order, err := h.storage.GetOrder(params.OrderUUID)
	if err != nil {
		return nil, err
	}

	// Преобразуем строковый PayOrderRequestPaymentMethod в PaymentMethodEnum для payment сервиса
	paymentMethodEnum := convertPaymentMethodToEnum(method)
	payOrderMessage := paymentV1.PayOrderMessage{
		OrderUuid:     params.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentMethodEnum,
	}
	transactionUuid, err := payOrderCall(&paymentV1.PayOrderRequest{PayOrderMessage: &payOrderMessage})
	if err != nil {
		// КРИТИЧНО: Заказ уже в статусе PAID, но транзакции нет
		// Нужно либо откатить статус, либо вернуть ошибку
		// Пока возвращаем ошибку, чтобы клиент знал о проблеме
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	order_pay_resp := &orderV1.PayOrderResponse{
		//		TransactionUUID: "7b5c38cb-57b9-4f2e-9a0a-c518add9ccaa",
		TransactionUUID: transactionUuid,
		// тут должен наверное генерироваться uuid транзакции
	}
	return order_pay_resp, nil
}

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order_uuid := &orderV1.CancelOrderParams{
		OrderUUID: params.OrderUUID,
	}

	//	h.storage.CancelOrderByUUID(order_uuid)

	//	return &orderV1.CancelOrderNoContent{}, nil
	return h.storage.CancelOrderByUUID(order_uuid)
}

// NewError создает новую ошибку в формате GenericError
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {

	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Ошибка создания сервера OpenApi: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	// Написать кастомные мидлвари

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
