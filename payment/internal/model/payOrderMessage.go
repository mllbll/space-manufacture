package model

type PayOrderMessage struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethodEnum
}

type PaymentMethodEnum int32

const (
	PAYMENT_METHOD_ENUM_UNSPECIFIED    PaymentMethodEnum = 0
	PAYMENT_METHOD_ENUM_CARD           PaymentMethodEnum = 1
	PAYMENT_METHOD_ENUM_SBP            PaymentMethodEnum = 2
	PAYMENT_METHOD_ENUM_CREDIT_CARD    PaymentMethodEnum = 3
	PAYMENT_METHOD_ENUM_INVESTOR_MONEY PaymentMethodEnum = 4
)

type PayOrderResponse struct {
	TransactionUUID string
}

type PayOrderRequest struct {
	PayOrderMessage PayOrderMessage
}
