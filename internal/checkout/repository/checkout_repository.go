package repository

import (
	"ps-eniqilo-store/internal/checkout/dto"
)

type CheckoutRepository interface {
	GetCheckoutHistory(params map[string]interface{}) ([]dto.CheckOutHistoryResp, error)
}
