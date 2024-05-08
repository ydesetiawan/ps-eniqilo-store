package service

import (
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/repository"
	"ps-eniqilo-store/pkg/errs"
	"ps-eniqilo-store/pkg/helper"
)

type ProductService interface {
	CreateProduct(*dto.ProductReq) (dto.ProductResp, error)
	UpdateProduct(*dto.ProductReq, int64) (dto.ProductResp, error)
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

func (s *productService) UpdateProduct(request *dto.ProductReq, productId int64) (dto.ProductResp, error) {
	_, err := s.productRepository.GetProductByID(productId)
	if err != nil {
		return dto.ProductResp{}, errs.NewErrDataNotFound("cat id is not found", productId, errs.ErrorData{})
	}

	savedProduct, err := s.productRepository.UpdateProduct(request, productId)
	if err != nil {
		return dto.ProductResp{}, err
	}

	return dto.ProductResp{
		ID:        helper.IntToString(savedProduct.ID),
		CreatedAt: savedProduct.CreatedAtFormatter,
	}, nil
}
