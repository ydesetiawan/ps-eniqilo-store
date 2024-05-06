package service

import (
	"ps-cats-social/internal/product/model"
	"ps-cats-social/internal/product/repository"
)

type ProductServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewProductServiceImpl(productRepository repository.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{productRepository: productRepository}
}

func (s *ProductServiceImpl) CreateProduct(product *model.Product) error {
	//TODO implement me
	panic("implement me")
}
