package service

import "ps-cats-social/internal/product/model"

type ProductService interface {
	CreateProduct(*model.Product) error
}
