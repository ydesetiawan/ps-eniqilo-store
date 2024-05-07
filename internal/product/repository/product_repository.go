package repository

import (
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/model"
)

type ProductRepository interface {
	GetProductByID(int64) (model.Product, error)
	CreateProduct(*dto.ProductReq) (model.Product, error)
	UpdateProduct(*dto.ProductReq, int64) (model.Product, error)
	DeleteProduct(int64) error
}
