package repository

import (
	"github.com/jmoiron/sqlx"
	"ps-eniqilo-store/internal/checkout/model"
)

type checkoutRepository struct {
	db *sqlx.DB
}

func NewCheckoutRepositoryImpl(db *sqlx.DB) CheckoutRepository {
	return &checkoutRepository{db: db}
}

func (c checkoutRepository) GetCheckoutHistory(params map[string]interface{}) ([]model.Checkout, error) {
	//TODO implement me
	panic("implement me")
}
