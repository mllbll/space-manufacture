package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/mllbll/space-manufacture/order/internal/api/order/v1"
	paymentClientV1 "github.com/mllbll/space-manufacture/order/internal/client/grpc/payment/v1"
	inventoryClientV1 "github.com/mllbll/space-manufacture/order/internal/client/grpc/inventory/v1"
	orderRepository "github.com/mllbll/space-manufacture/order/internal/repository/order"
	orderService "github.com/mllbll/space-manufacture/order/internal/service/order"
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

func main() {
	// Создаем gRPC соединения
	paymentConn, err := grpc.NewClient(
		paymentServiceServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}
	defer paymentConn.Close()

	inventoryConn, err := grpc.NewClient(
		inventoryServiceServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to inventory service: %v", err)
	}
	defer inventoryConn.Close()

	// Создаем клиенты
	paymentGeneratedClient := paymentV1.NewPaymentServiceClient(paymentConn)
	paymentClient := paymentClientV1.NewClient(paymentGeneratedClient)

	inventoryGeneratedClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	inventoryClient := inventoryClientV1.NewClient(inventoryGeneratedClient)

	// Создаем слои приложения
	repo := orderRepository.NewRepository()
	service := orderService.NewService(repo, paymentClient, inventoryClient)
	api := orderAPI.NewAPI(service)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatalf("Ошибка создания сервера OpenAPI: %v", err)
	}

	// Настраиваем HTTP роутер
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	// Настраиваем HTTP сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запускаем сервер в горутине
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

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
