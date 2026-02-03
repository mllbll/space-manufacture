package order

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	def "github.com/mllbll/space-manufacture/order/internal/repository"

	"github.com/mllbll/space-manufacture/order/internal/client/db"

	repoModel "github.com/mllbll/space-manufacture/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type postgresRepository struct {
	mu sync.RWMutex
	pool *pgxpool.Pool
}

func NewPostgresRepository(database *db.DB) *postgresRepository {
	return &postgresRepository{
		pool: database.GetPool(),
	}
}
