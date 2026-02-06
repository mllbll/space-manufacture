package order

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mllbll/space-manufacture/order/internal/model"
	repoConverter "github.com/mllbll/space-manufacture/order/internal/repository/converter"
)

func (r *postgresRepository) Create(ctx context.Context, req model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoReq := repoConverter.OrderToRepoModel(req)

	createBuilder := sq.Insert("orders").PlaceholderFormat(sq.Dollar).Columns("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status").Values(repoReq.OrderUUID, repoReq.UserUUID, repoReq.PartUUIDs, repoReq.TotalPrice, repoReq.TransactionUUID, repoReq.PaymentMethod, repoReq.Status)

	query, args, err := createBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Print("Заказ с таким uuid уже существует")
        return model.ErrOrderConflict
    }
		return err
	}

	return nil
}
