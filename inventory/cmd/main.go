package main

import (
	"fmt"
	"github.com/mllbll/space-manufacture/inventory/internal/interceptor"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/mllbll/space-manufacture/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/mllbll/space-manufacture/inventory/internal/repository/part"
	inventoryService "github.com/mllbll/space-manufacture/inventory/internal/service/part"
	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	// регистрируем сервис
	repo := inventoryRepository.NewRepository()
	// Инициализируем тестовые данные
	repo.InitTestData()
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC recovered in gRPC server: %v\n", r)
				// Передаем панику дальше для полного стека
				panic(r)
			}
		}()
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server ...")
	s.GracefulStop()
	log.Println("Server Stopped")

}
