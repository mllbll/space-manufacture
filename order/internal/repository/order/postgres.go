package order

import (
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	def "github.com/mllbll/space-manufacture/order/internal/repository"

)

var _ def.OrderRepository = (*postgresRepository)(nil)

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type postgresRepository struct {
	mu sync.RWMutex
	pool *pgxpool.Pool
}

func NewPostgresRepository(dbPool *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		pool: dbPool,
	}
}
