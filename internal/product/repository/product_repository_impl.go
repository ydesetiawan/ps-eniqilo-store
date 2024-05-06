package repository

import (
	"github.com/jmoiron/sqlx"
	"ps-eniqilo-store/internal/product/model"
)

type ProductRepositoryImpl struct {
	db *sqlx.DB
}

func NewProductRepositoryImpl(db *sqlx.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) CreateProduct(product *model.Product) error {
	//TODO implement me
	panic("implement me")
}
