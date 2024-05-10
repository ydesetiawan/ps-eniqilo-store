package repository

import (
	"ps-eniqilo-store/internal/checkout/dto"
	"ps-eniqilo-store/internal/checkout/model"

	"github.com/jmoiron/sqlx"
)

type CheckoutRepository interface {
	GetCheckoutHistory(params map[string]interface{}) ([]dto.CheckOutHistoryResp, error)
	CreateCheckout(*sqlx.Tx, *model.Checkout) (int64, error)
	CreateCheckoutDetail(*sqlx.Tx, *model.CheckoutDetail) error
}
