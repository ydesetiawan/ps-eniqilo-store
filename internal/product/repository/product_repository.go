package repository

import (
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/model"
)

type ProductRepository interface {
	CreateProduct(*dto.ProductReq) (model.Product, error)
}
