package app

import (
	"context"
	"fmt"

	"github.com/mllbll/space-manufacture/inventory/internal/config"
	"github.com/mllbll/space-manufacture/inventory/internal/repository"
	"github.com/mllbll/space-manufacture/inventory/internal/service"
	"platform/pkg/closer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	

	inventoryV1API "github.com/mllbll/space-manufacture/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/mllbll/space-manufacture/inventory/internal/repository/part"
	inventoryService "github.com/mllbll/space-manufacture/inventory/internal/service/part"
	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1Api inventoryV1.InventoryServiceServer

	inventoryService service.InventoryService

	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client

	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1Api == nil {
		d.inventoryV1Api = inventoryV1API.NewAPI(d.PartService(ctx))
	}
	
	return d.inventoryV1Api

}


func (d *diContainer) PartService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.PartRepository(ctx))
	}

	return d.inventoryService

}

func (d *diContainer) PartRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewMongoCollection(d.MongoDBHandle(ctx))
	}

	return d.inventoryRepository

}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to Mongo: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
