package order

import (
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	def "github.com/mllbll/space-manufacture/order/internal/repository"

	"github.com/mllbll/space-manufacture/order/internal/client/db"
)

var _ def.OrderRepository = (*postgresRepository)(nil)

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type postgresRepository struct {
	mu sync.RWMutex
	pool *pgxpool.Pool
}

func NewPostgresRepository(database *db.DB) *postgresRepository {
	return &postgresRepository{
		pool: database.GetPool(),
	}
}
