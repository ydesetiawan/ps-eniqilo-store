package repository

import (
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/model"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	GetProductByID(int64) (model.Product, error)
	GetProductByIDs([]int64) ([]model.Product, error)
	CreateProduct(*dto.ProductReq) (model.Product, error)
	GetProducts(params map[string]interface{}) ([]model.Product, error)
	UpdateProduct(*dto.ProductReq, int64) (model.Product, error)
	DeleteProduct(int64) error
	DecreaseStock(*sqlx.Tx, int64, int) error
}
