package model

type Order struct {
	// Уникальный идентификатор заказа (UUID)
	OrderUUID string

	//UUID пользователя
	UserUUID string

	//Список UUID деталей
	PartUUIDs []string

	//Итоговая стоимость
	TotalPrice float32

	//uuid транзакции (если оплачен)
	TransactionUUID *string

	//способ оплаты
	PaymentMethod *int32

	//статус оплаты (если оплачен)
	Status string
}


type CreateOrderRequest struct {
	UserUUID string

	PartUUIDs []string
}

type CreateOrderResponse struct {
	OrderUUID string

	TotalPrice float32
}

type GetOrderResponce struct {
	// Уникальный идентификатор заказа (UUID)
	OrderUUID string

	//UUID пользователя
	UserUUID string

	//Список UUID деталей
	PartUUIDs []string

	//Итоговая стоимость
	TotalPrice float32

	//uuid транзакции (если оплачен)
	TransactionUUID *string

	//способ оплаты
	PaymentMethod *int32

	//статус оплаты
	Status string
}

type PayOrderRequest struct {
	PaymentMethod int32
}

type PayOrderResponse struct {
	TransactionUUID string
}
