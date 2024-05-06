package repository

import "ps-eniqilo-store/internal/product/model"

type ProductRepository interface {
	CreateProduct(*model.Product) error
}
