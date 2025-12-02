package model

import (
	"errors"
	// orderV1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
)

var ErrOrderNotFound = errors.New("order not found")

// var ErrOrderConflict = orderV1.CancelOrderConflict{}
var ErrOrderConflict = errors.New("conflict")
var ErrOrderNoContent = errors.New("No content")
