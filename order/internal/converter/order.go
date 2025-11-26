package converter

import (
	"github.com/mllbll/space-manufacture/order/internal/model"
	orderV1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
	"github.com/samber/lo"
)

func OrderToModel(req *orderV1.Order) model.Order {
	var transactionUUID *string
	if req.TransactionUUID.Set {
		transactionUUID = lo.ToPtr(req.TransactionUUID.Value)
	}

	var paymentMethod *model.PaymentMethodEnum
	if req.PaymentMethod.Set {
		paymentMethod = lo.ToPtr(model.PaymentMethodEnum(req.PaymentMethod.Value))
	}
	return model.Order{
		OrderUUID: req.OrderUUID,
		UserUUID: req.UserUUID,
		PartUUIDs: req.PartUuids,
		TotalPrice: req.TotalPrice,
		TransactionUUID: transactionUUID,
		Status: string(req.Status),
		PaymentMethod: paymentMethod,
	}
}

func OrderToOpenAPI(req model.Order) *orderV1.Order {

	var transactionUUID orderV1.OptString
	if req.TransactionUUID != nil {
		transactionUUID = orderV1.NewOptString(*req.TransactionUUID)
	}

	var paymentMethod orderV1.OptOrderPaymentMethod
	if req.PaymentMethod != nil {
		paymentMethod = orderV1.NewOptOrderPaymentMethod(orderV1.OrderPaymentMethod(*req.PaymentMethod))
	}
	return &orderV1.Order{
		OrderUUID: req.OrderUUID,
		UserUUID: req.UserUUID,
		PartUuids: req.PartUUIDs,
		TotalPrice: float32(req.TotalPrice),
		TransactionUUID: transactionUUID,
		Status: orderV1.OrderStatus(req.Status),
		PaymentMethod: paymentMethod,
	}
}

func CreateOrderRequestToModel(req *orderV1.CreateOrderRequest) model.CreateOrderRequest {
	return model.CreateOrderRequest{
		UserUUID: req.UserUUID,
		PartUUIDs: req.PartUuids,
	}
}

func CreateOrderRequestToOpenAPI(req model.CreateOrderRequest) *orderV1.CreateOrderRequest {
	return &orderV1.CreateOrderRequest{
		UserUUID: req.UserUUID,
		PartUuids: req.PartUUIDs,
	}
}

func CreateOrderResponseToModel(req *orderV1.CreateOrderResponse) model.CreateOrderResponse {
	return model.CreateOrderResponse{
		OrderUUID: req.OrderUUID.Value,
		TotalPrice: req.TotalPrice.Value,
	}
}

func CreateOrderResponseToOpenAPI(req model.CreateOrderResponse) *orderV1.CreateOrderResponse {
	return &orderV1.CreateOrderResponse{
		OrderUUID: req.OrderUUID,
		TotalPrice: req.TotalPrice,
	}
}

func PayOrderRequestToModel(req *orderV1.PayOrderRequest) model.PayOrderRequest {
	// Конвертируем PayOrderRequestPaymentMethod (string) в model.PaymentMethodEnum (int32)
	var paymentMethod model.PaymentMethodEnum
	switch req.PaymentMethod {
	case orderV1.NewOptPayOrderRequestPaymentMethod(orderV1.PayOrderRequestPaymentMethod0):
		paymentMethod = model.PAYMENT_METHOD_ENUM_UNSPECIFIED
	case orderV1.NewOptPayOrderRequestPaymentMethod(orderV1.PayOrderRequestPaymentMethod1):
		paymentMethod = model.PAYMENT_METHOD_ENUM_CARD
	case orderV1.NewOptPayOrderRequestPaymentMethod(orderV1.PayOrderRequestPaymentMethod2):
		paymentMethod = model.PAYMENT_METHOD_ENUM_SBP
	case orderV1.NewOptPayOrderRequestPaymentMethod(orderV1.PayOrderRequestPaymentMethod3):
		paymentMethod = model.PAYMENT_METHOD_ENUM_CREDIT_CARD
	case orderV1.NewOptPayOrderRequestPaymentMethod(orderV1.PayOrderRequestPaymentMethod4):
		paymentMethod = model.PAYMENT_METHOD_ENUM_INVESTOR_MONEY
	default:
		paymentMethod = model.PAYMENT_METHOD_ENUM_UNSPECIFIED
	}

	return model.PayOrderRequest{
		PaymentMethod: paymentMethod,
	}
}

func PayOrderRequestToOpenAPI(req model.PayOrderRequest) *orderV1.PayOrderRequest {
	// Конвертируем model.PaymentMethodEnum (int32) в PayOrderRequestPaymentMethod (string)
	var paymentMethod orderV1.PayOrderRequestPaymentMethod
	switch req.PaymentMethod {
	case model.PAYMENT_METHOD_ENUM_UNSPECIFIED:
		paymentMethod = orderV1.PayOrderRequestPaymentMethod0
	case model.PAYMENT_METHOD_ENUM_CARD:
		paymentMethod = orderV1.PayOrderRequestPaymentMethod1
	case model.PAYMENT_METHOD_ENUM_SBP:
		paymentMethod = orderV1.PayOrderRequestPaymentMethod2
	case model.PAYMENT_METHOD_ENUM_CREDIT_CARD:
		paymentMethod = orderV1.PayOrderRequestPaymentMethod3
	case model.PAYMENT_METHOD_ENUM_INVESTOR_MONEY:
		paymentMethod = orderV1.PayOrderRequestPaymentMethod4
	default:
		paymentMethod = orderV1.PayOrderRequestPaymentMethod4
	}

	return &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.NewOptPayOrderRequestPaymentMethod(paymentMethod),
	}
}

func GetOrderResponseToModel(req *orderV1.Order) model.GetOrderResponce {
	var transactionUUID *string
	if req.TransactionUUID.Set {
		transactionUUID = lo.ToPtr(req.TransactionUUID.Value)
	}

	var paymentMethod *model.PaymentMethodEnum
	if req.PaymentMethod.Set {
		paymentMethod = lo.ToPtr(model.PaymentMethodEnum(req.PaymentMethod.Value))
	}
	return model.GetOrderResponce{
		OrderUUID: req.OrderUUID,
		UserUUID: req.UserUUID,
		PartUUIDs: req.PartUuids,
		TotalPrice: req.TotalPrice,
		TransactionUUID: transactionUUID,
		Status: string(req.Status),
		PaymentMethod: paymentMethod,
	}
}


func GetOrderResponseToOpenAPI(req model.GetOrderResponce) *orderV1.Order {

	var transactionUUID orderV1.OptString
	if req.TransactionUUID != nil {
		transactionUUID = orderV1.NewOptString(*req.TransactionUUID)
	}

	var paymentMethod orderV1.OptOrderPaymentMethod
	if req.PaymentMethod != nil {
		paymentMethod = orderV1.NewOptOrderPaymentMethod(orderV1.OrderPaymentMethod(*req.PaymentMethod))
	}
	return &orderV1.Order{
		OrderUUID: req.OrderUUID,
		UserUUID: req.UserUUID,
		PartUuids: req.PartUUIDs,
		TotalPrice: float32(req.TotalPrice),
		TransactionUUID: transactionUUID,
		Status: orderV1.OrderStatus(req.Status),
		PaymentMethod: paymentMethod,
	}
}

func PayOrderResponseToModel(req *orderV1.PayOrderResponse) model.PayOrderResponse {
	return model.PayOrderResponse{
		TransactionUUID: req.TransactionUUID,
	}
}

func PayOrderResponseToOpenAPI(req model.PayOrderResponse) *orderV1.PayOrderResponse {
	return &orderV1.PayOrderResponse{
		TransactionUUID: req.TransactionUUID,
	}
}


