package repository

import (
	"github.com/jmoiron/sqlx"
	"ps-eniqilo-store/internal/product/model"
)

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepositoryImpl(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *model.Product) error {
	//TODO implement me
	panic("implement me")
}
