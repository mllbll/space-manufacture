package db

import (
	"context"
	"fmt"
	"log"

	"github.com/mllbll/space-manufacture/inventory/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
	db     *mongo.Database
}

const configPath = "./../deploy/compose/inventory/.env"

func NewDb() (*DB, error) {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to mongo database: %v\n", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("MongoDB не доступна, ошибка: %v\n", err)
		return nil, err
	}

	database := client.Database("inventory")
	return &DB{
		client: client,
		db:     database,
	}, nil

}

func (d *DB) Close() error {
	return d.client.Disconnect(context.Background())
}

func (d *DB) GetDataBase() *mongo.Database {
	return d.db
}

func (d *DB) GetClient() *mongo.Client {
	return d.client
}
