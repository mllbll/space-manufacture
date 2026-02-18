package app

import (
	"context"
	"log"

	"platform/pkg/closer"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/mllbll/space-manufacture/order/internal/api/order/v1"
	grpcClient "github.com/mllbll/space-manufacture/order/internal/client/grpc"
	inventoryClientV1 "github.com/mllbll/space-manufacture/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/mllbll/space-manufacture/order/internal/client/grpc/payment/v1"
	"github.com/mllbll/space-manufacture/order/internal/config"
	"github.com/mllbll/space-manufacture/order/internal/migrator"
	"github.com/mllbll/space-manufacture/order/internal/repository"
	orderRepository "github.com/mllbll/space-manufacture/order/internal/repository/order"
	"github.com/mllbll/space-manufacture/order/internal/service"
	orderService "github.com/mllbll/space-manufacture/order/internal/service/order"
	orderV1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API orderV1.Handler

	orderService service.OrderService

	orderRepository repository.OrderRepository

	postgresClient *pgxpool.Pool

	inventoryClient grpc.ClientConnInterface

	inventoryClientHandler grpcClient.InventoryClient

	paymentClient grpc.ClientConnInterface

	paymentClientHandler grpcClient.PaymentClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PostgresClient(ctx context.Context) *pgxpool.Pool {
	if d.postgresClient == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			log.Printf("Failed to connect to database %v\n", err)
			return nil
		}

		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			log.Printf("failed to ping database %w", err)
			return nil
		}

		migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(pool), config.AppConfig().Postgres.MigrationDir())

		err = migratorRunner.Up()
		if err != nil {
			log.Fatalf("Ошибка миграции базы данных", err)
		}

		d.postgresClient = pool
	}
	return d.postgresClient
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewPostgresRepository(d.PostgresClient(ctx))
	}
	return d.orderRepository
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(d.OrderRepository(ctx), d.PaymentClientHandle(ctx), d.InventoryClientHandle(ctx))
	}

	return d.orderService

}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}
	return d.orderV1API
}

func (d *diContainer) PaymentClientHandle(ctx context.Context) grpcClient.PaymentClient {
	if d.paymentClientHandler == nil {
		conn := d.PaymentClient(ctx)
		paymentGeneratedClient := paymentV1.NewPaymentServiceClient(conn)
		d.paymentClientHandler = paymentClientV1.NewClient(paymentGeneratedClient)
	}

	return d.paymentClientHandler
}

func (d *diContainer) PaymentClient(ctx context.Context) grpc.ClientConnInterface {
	if d.paymentClient == nil {
		paymentConn, err := grpc.NewClient(config.AppConfig().PaymentGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal("Ошибка создания коннекта к payment сервису")
			return nil
		}

		closer.AddNamed("PaymentClient connection", func(ctx context.Context) error {
			return paymentConn.Close()
		})

		d.paymentClient = paymentConn
	}

	return d.paymentClient
}

func (d *diContainer) InventoryClient(ctx context.Context) grpc.ClientConnInterface {
	if d.inventoryClient == nil {
		inventoryConn, err := grpc.NewClient(config.AppConfig().InventoryGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal("Ошибка создания коннекта к inventory сервису")
			return nil
		}

		closer.AddNamed("InventoryClient Connection", func(ctx context.Context) error {
			return inventoryConn.Close()
		})

		d.inventoryClient = inventoryConn
	}

	return d.inventoryClient
}

func (d *diContainer) InventoryClientHandle(ctx context.Context) grpcClient.InventoryClient {
	if d.inventoryClientHandler == nil {
		conn := d.InventoryClient(ctx)
		inventoryGeneratedClient := inventoryV1.NewInventoryServiceClient(conn)
		d.inventoryClientHandler = inventoryClientV1.NewClient(inventoryGeneratedClient)
	}

	return d.inventoryClientHandler
}
