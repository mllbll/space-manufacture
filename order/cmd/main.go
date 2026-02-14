package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	"github.com/mllbll/space-manufacture/order/internal/client/db"
	inventoryClientV1 "github.com/mllbll/space-manufacture/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/mllbll/space-manufacture/order/internal/client/grpc/payment/v1"
	"github.com/mllbll/space-manufacture/order/internal/config"
	orderRepository "github.com/mllbll/space-manufacture/order/internal/repository/order"
	orderService "github.com/mllbll/space-manufacture/order/internal/service/order"
	orderV1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

const configPath = "./../deploy/compose/order/.env"

const (
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// const (
// 	httpPort = "8080"
// 	// Таймауты для HTTP-сервера
// 	readHeaderTimeout = 5 * time.Second
// 	shutdownTimeout   = 10 * time.Second
// 	// Порты gRPC сервисов
// 	paymentServiceServerAddress   = "localhost:50051"
// 	inventoryServiceServerAddress = "localhost:50052"
// )

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config:%w", err))
	}

	// Создаем gRPC соединения
	paymentConn, err := grpc.NewClient(
		config.AppConfig().PaymentGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}
	defer paymentConn.Close()

	// Коннект к inventory сервису
	inventoryConn, err := grpc.NewClient(
		config.AppConfig().InventoryGRPC.Address(),
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

	// Открываем коннект с БД
	database, err := db.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Закрываем коннект что бы не висел
	defer database.Close()

	// Создаем слои приложения
	repo := orderRepository.NewPostgresRepository(database)
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
		Addr:              config.AppConfig().OrderHTTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", config.AppConfig().OrderHTTTP.Address())
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
