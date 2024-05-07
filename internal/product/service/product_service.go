package service

import (
	"ps-eniqilo-store/internal/product/model"
	"ps-eniqilo-store/internal/product/repository"
)

type ProductService interface {
	CreateProduct(*model.Product) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductServiceImpl(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) CreateProduct(product *model.Product) error {
	//TODO implement me
	panic("implement me")
}
