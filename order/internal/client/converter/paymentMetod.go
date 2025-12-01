package converter

import (
	"github.com/mllbll/space-manufacture/order/internal/model"
	generatedPaymentV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/payment/v1"
)
// PaymentMethodEnumToString конвертирует model.PaymentMethodEnum (int32) в строку
// для использования в payment client
func PaymentMethodEnumToString(method model.PaymentMethodEnum) string {
	// Используем маппинг из proto для конвертации int32 в строку
	if str, ok := generatedPaymentV1.PaymentMethodEnum_name[int32(method)]; ok {
		return str
	}
	// Возвращаем UNSPECIFIED по умолчанию, если значение не найдено
	return generatedPaymentV1.PaymentMethodEnum_name[int32(model.PAYMENT_METHOD_ENUM_UNSPECIFIED)]
}


