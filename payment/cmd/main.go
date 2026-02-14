package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mllbll/space-manufacture/payment/internal/config"
	"github.com/mllbll/space-manufacture/payment/internal/interceptor"

	paymentV1API "github.com/mllbll/space-manufacture/payment/internal/api/payment/v1"
	paymentRepository "github.com/mllbll/space-manufacture/payment/internal/repository/payment"
	paymentService "github.com/mllbll/space-manufacture/payment/internal/service/payment"

	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

const configPath = "./../deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
	panic(fmt.Errorf("failed to load config: %w", err))
	}

	lis, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
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
		log.Printf("gRPC server listening on %s\n", config.AppConfig().PaymentGRPC.Address())
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
