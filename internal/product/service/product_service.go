package service

import (
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/repository"
	"ps-eniqilo-store/pkg/helper"
)

type ProductService interface {
	CreateProduct(*dto.ProductReq) (dto.ProductResp, error)
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductServiceImpl(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) CreateProduct(product *dto.ProductReq) (dto.ProductResp, error) {
	savedProduct, err := s.productRepository.CreateProduct(product)
	if err != nil {
		return dto.ProductResp{}, err
	}

	return dto.ProductResp{
		ID:        helper.IntToString(savedProduct.ID),
		CreatedAt: savedProduct.CreatedAtFormatter,
	}, nil
}
