package model

import ()

type Order struct {
	// Уникальный идентификатор заказа (UUID)
	OrderUUID string

	//UUID пользователя
	UserUUID string

	//Список UUID деталей
	PartUUIDs []string

	//Итоговая стоимость
	TotalPrice float64

	//uuid транзакции (если оплачен)
	TransactionUUID *string

	//статус оплаты (если оплачен)
	status *string

	//способ оплаты
	PaymentMethod PaymentMethodEnum
}

// Enum структура для методов оплаты
type PaymentMethodEnum int32

const (
	//Неизвестный способ
	PAYMENT_METHOD_ENUM_UNSPECIFIED    PaymentMethodEnum = 0
	// Банковская карта
	PAYMENT_METHOD_ENUM_CARD           PaymentMethodEnum = 1
	//Система быстрых платежей
	PAYMENT_METHOD_ENUM_SBP            PaymentMethodEnum = 2
	//Кредитная карта
	PAYMENT_METHOD_ENUM_CREDIT_CARD    PaymentMethodEnum = 3
	//Деньги инвестора
	PAYMENT_METHOD_ENUM_INVESTOR_MONEY PaymentMethodEnum = 4
)

type CreateOrderRequest struct {
	UserUUID string

	PartUUIDs []string
}


type CreateOrderResponse struct {
	UserUUID string

	PartUUIDs []string
}

type GetOrderResponce struct {
	Order Order
}


type PayOrderRequest struct {
	UserUUID string

	PartUUIDs []string
}

type PayOrderResponse struct {
	TransactionUUID string
}
