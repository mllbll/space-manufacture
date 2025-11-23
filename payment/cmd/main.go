package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"payment/internal/interceptor"

	paymentV1API "github.com/mllbll/space-manufacture/payment/internal/payment/v1"
	paymentRepository "github.com/space-manufacture/payment/internal/repository/payment"
	paymentService "github.com/space-manufacture/payment/internal/service/payment"

	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

const grpcPort = 50051

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

	// создаем grpc сервер с интерсептором logger
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	// регистрируем сервис
	repo := paymentRepository.NewRepository()
	service := paymentService.NewService(repo)
	api := paymentV1API.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	// включаем рефлексию
	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down gRPC server ...")
	s.GracefulStop()
	log.Printf("Server Stopped")

}
