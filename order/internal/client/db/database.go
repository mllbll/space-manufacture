package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mllbll/space-manufacture/order/internal/config"
	"github.com/mllbll/space-manufacture/order/internal/migrator"
)

type DB struct {
	pool *pgxpool.Pool
}

const configPath = "./../deploy/compose/order/.env"

func NewDB() (*DB, error) {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Printf("failed to load .env file: %v\n", err)
	// 	return nil, err
	// }
	//
	// dbURI := os.Getenv("DB_URI")

	pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		log.Printf("Failed to connect to database %v\n", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database %w", err)
	}

	migrationDir := config.AppConfig().Postgres.MigrationDir()
	migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(pool), migrationDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции бд %v\n", err)
		return nil, err
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *DB) GetPool() *pgxpool.Pool {
	return db.pool
}
