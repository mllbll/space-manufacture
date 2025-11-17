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

	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

const grpcPort = 50051

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
	// мапка со значениями transaction_uuid:payOrderMessage
	mu               sync.RWMutex
	payOrderMessages map[string]*paymentV1.PayOrderMessage
}

func (s *paymentService) Pay(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.PayOrderMessage == nil {
		return nil, fmt.Errorf("pay_order_message is required")
	}

	newTransactionUUID := uuid.New().String()

	newPayOrderMessage := &paymentV1.PayOrderMessage{
		OrderUuid:     req.PayOrderMessage.OrderUuid,
		UserUuid:      req.PayOrderMessage.UserUuid,
		PaymentMethod: req.PayOrderMessage.PaymentMethod,
	}

	s.payOrderMessages[newTransactionUUID] = newPayOrderMessage

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", newTransactionUUID)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: newTransactionUUID,
	}, nil

}

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

	s := grpc.NewServer()

	service := &paymentService{
		payOrderMessages: make(map[string]*paymentV1.PayOrderMessage),
	}

	paymentV1.RegisterPaymentServiceServer(s, service)

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
