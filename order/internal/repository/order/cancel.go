package order

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/mllbll/space-manufacture/order/internal/model"
	repoConverter "github.com/mllbll/space-manufacture/order/internal/repository/converter"
	repoModel "github.com/mllbll/space-manufacture/order/internal/repository/model"
)

func (r *postgresRepository) Cancel(ctx context.Context, param string) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cancelBuilder := sq.Update("orders").PlaceholderFormat(sq.Dollar).Set("status", "CANCELLED").Where(sq.Eq{"order_uuid": param}).Suffix("RETURNING order_uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status")

	query, args, err := cancelBuilder.ToSql()
	if err != nil {
		return model.Order{}, err
	}
	
	var row repoModel.Order
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&row.OrderUUID,
		&row.UserUUID,
		&row.PartUUIDs,
		&row.TotalPrice,
		&row.TransactionUUID,
		&row.PaymentMethod,
		&row.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		log.Printf("failed to cancel order: %v\n", err)
		return model.Order{}, err
	}

	return repoConverter.OrderToModel(row), nil
}
