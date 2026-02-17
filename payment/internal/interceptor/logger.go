package interceptor

import (
	"context"
	"fmt"
	"log"
	"path"
	"platform/pkg/logger"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Извлекаем имя метода из полного пути
		method := path.Base(info.FullMethod)

		// Логируем начало вызова метода
		log.Printf("🚀 Started gRPC method %s\n", method)
		logger.Info(ctx, "🚀 Started gRPC method", zap.String(method, ""))

		// Засекаем время начала выполнения
		startTime := time.Now()

		// Вызываем обработчик
		resp, err := handler(ctx, req)

		// Вычисляем длительность выполнения
		duration := time.Since(startTime)

		// Форматируем сообщение в зависимости от результата
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error(ctx, fmt.Sprintf("❌ Finished gRPC method %s with code %s: %v (took: %v)\n", method, st.Code(), err, duration), zap.Error(err))
			log.Printf("❌ Finished gRPC method %s with code %s: %v (took: %v)\n", method, st.Code(), err, duration)
		} else {
			log.Printf("✅ Finished gRPC method %s successfully (took: %v)\n", method, duration)
			logger.Info(ctx, fmt.Sprintf("✅ Finished gRPC method %s successfully (took: %v)\n", method, duration))
		}

		return resp, err
	}
}

func LoggerInterceptor2() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Логируем начало вызова метода
		log.Printf("🚀 Bla\n")

		resp, err := handler(ctx, req)

		return resp, err
	}
}
