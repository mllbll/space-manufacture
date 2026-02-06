package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/mllbll/space-manufacture/order/internal/model"
	repoConverter "github.com/mllbll/space-manufacture/order/internal/repository/converter"
	repoModel "github.com/mllbll/space-manufacture/order/internal/repository/model"
)

func (r *postgresRepository) Get(ctx context.Context, params string) (model.GetOrderResponce, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	getBuilder := sq.Select("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status").From("orders").Where(sq.Eq{"order_uuid": params}).PlaceholderFormat(sq.Dollar)

	query, args, err := getBuilder.ToSql()

	if err != nil {
		log.Print("failed to query sql")
		return model.GetOrderResponce{}, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("failed to select orders : %v\n", err)
		return model.GetOrderResponce{}, err
	}

	defer rows.Close()

	var res repoModel.GetOrderResponce

	if !rows.Next() {
		return model.GetOrderResponce{}, model.ErrOrderNotFound
	}

	err = rows.Scan(&res.OrderUUID, &res.UserUUID, &res.PartUUIDs, &res.TotalPrice, &res.TransactionUUID, &res.PaymentMethod, &res.Status)

	if err != nil {
		log.Printf("failed to scan order with uuid: %v\n", err)
		return model.GetOrderResponce{}, err
	}

	return repoConverter.GetOrderResponseToModel(res), nil
}
