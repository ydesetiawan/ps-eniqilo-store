package repository

import (
	"ps-eniqilo-store/internal/checkout/model"
)

type CheckoutRepository interface {
	GetCheckoutHistory(params map[string]interface{}) ([]model.Checkout, error)
}
