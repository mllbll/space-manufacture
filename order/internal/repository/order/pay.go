package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/mllbll/space-manufacture/order/internal/model"
)

func (r *postgresRepository) Pay(ctx context.Context, param string, req model.PayOrderRequest, transactionUUID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	payBuilder := sq.Update("orders").PlaceholderFormat(sq.Dollar).Set("status", "PAID").Set("transaction_uuid", transactionUUID).Set("payment_method", req.PaymentMethod).Where(sq.Eq{"order_uuid": param})

	query, args, err := payBuilder.ToSql()
	if err != nil {
		return err
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update order : %v\n", err)
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}
	return nil
}
