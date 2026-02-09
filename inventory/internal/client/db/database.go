package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewDb() (*DB, error) {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return nil, err
	}
	dbURI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
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
