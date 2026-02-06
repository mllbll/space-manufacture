package converter

import (
	"github.com/mllbll/space-manufacture/order/internal/model"
	repoModel "github.com/mllbll/space-manufacture/order/internal/repository/model"

)

func OrderToModel(req repoModel.Order) model.Order {
	return model.Order{
		OrderUUID:       req.OrderUUID,
		UserUUID:        req.UserUUID,
		PartUUIDs:       req.PartUUIDs,
		TotalPrice:      req.TotalPrice,
		TransactionUUID: req.TransactionUUID,
		Status:          req.Status,
		PaymentMethod:   Int32ToPaymentMethodEmum(req.PaymentMethod),
	}
}

func OrderToRepoModel(req model.Order) repoModel.Order {
	return repoModel.Order{
		OrderUUID:       req.OrderUUID,
		UserUUID:        req.UserUUID,
		PartUUIDs:       req.PartUUIDs,
		TotalPrice:      req.TotalPrice,
		TransactionUUID: req.TransactionUUID,
		Status:          req.Status,
		PaymentMethod:   PaymentMethodEnumToInt32(req.PaymentMethod),
	}
}

func Int32ToPaymentMethodEmum(req *int32) *model.PaymentMethodEnum {
	if req == nil {
		return nil
	}
	res := model.PaymentMethodEnum(*req)
	return &res
}

func PaymentMethodEnumToInt32(req *model.PaymentMethodEnum) *int32 {
	if req == nil {
		return nil
	}
	res := int32(*req)
	return &res
}

func CreateOrderRequestToModel(req repoModel.CreateOrderRequest) model.CreateOrderRequest {
	return model.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUUIDs,
	}
}

func CreateOrderRequestToRepoModel(req model.CreateOrderRequest) repoModel.CreateOrderRequest {
	return repoModel.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUUIDs,
	}
}

func CreateOrderResponseToModel(req repoModel.CreateOrderResponse) model.CreateOrderResponse {
	return model.CreateOrderResponse{
		OrderUUID:  req.OrderUUID,
		TotalPrice: req.TotalPrice,
	}
}

func CreateOrderResponseToRepoModel(req model.CreateOrderResponse) repoModel.CreateOrderResponse {
	return repoModel.CreateOrderResponse{
		OrderUUID:  req.OrderUUID,
		TotalPrice: req.TotalPrice,
	}
}

func PayOrderRequestToModel(req repoModel.PayOrderRequest) model.PayOrderRequest {
	return model.PayOrderRequest{
		PaymentMethod: *Int32ToPaymentMethodEmum(&req.PaymentMethod),
	}
}

func PayOrderRequestToRepoModel(req model.PayOrderRequest) repoModel.PayOrderRequest {
	return repoModel.PayOrderRequest{
		PaymentMethod: int32(req.PaymentMethod),
	}
}

func GetOrderResponseToModel(req repoModel.GetOrderResponce) model.GetOrderResponce {
	return model.GetOrderResponce{
		OrderUUID:       req.OrderUUID,
		UserUUID:        req.UserUUID,
		PartUUIDs:       req.PartUUIDs,
		TotalPrice:      req.TotalPrice,
		TransactionUUID: req.TransactionUUID,
		Status:          req.Status,
		PaymentMethod:   Int32ToPaymentMethodEmum(req.PaymentMethod),
	}
}

func GetOrderResponseToRepoModel(req model.GetOrderResponce) repoModel.GetOrderResponce {
	return repoModel.GetOrderResponce{
		OrderUUID:       req.OrderUUID,
		UserUUID:        req.UserUUID,
		PartUUIDs:       req.PartUUIDs,
		TotalPrice:      req.TotalPrice,
		TransactionUUID: req.TransactionUUID,
		Status:          req.Status,
		PaymentMethod:   PaymentMethodEnumToInt32(req.PaymentMethod),
	}
}

func PayOrderResponseToModel(req repoModel.PayOrderResponse) model.PayOrderResponse {
	return model.PayOrderResponse{
		TransactionUUID: req.TransactionUUID,
	}
}

func PayOrderResponseToRepoModel(req model.PayOrderResponse) repoModel.PayOrderResponse {
	return repoModel.PayOrderResponse{
		TransactionUUID: req.TransactionUUID,
	}
}
