package service

import (
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/model"
	"ps-eniqilo-store/internal/product/repository"
	"ps-eniqilo-store/pkg/errs"
	"ps-eniqilo-store/pkg/helper"
)

type ProductService interface {
	CreateProduct(*dto.ProductReq) (dto.ProductSuccessResp, error)
	GetProducts(map[string]interface{}) ([]model.Product, error)
	UpdateProduct(*dto.ProductReq, int64) (dto.ProductSuccessResp, error)
	DeleteProduct(int64) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductServiceImpl(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) CreateProduct(product *dto.ProductReq) (dto.ProductSuccessResp, error) {
	savedProduct, err := s.productRepository.CreateProduct(product)
	if err != nil {
		return dto.ProductSuccessResp{}, err
	}

	return dto.ProductSuccessResp{
		ID:        helper.IntToString(savedProduct.ID),
		CreatedAt: savedProduct.CreatedAtFormatter,
	}, nil
}

func (s *productService) GetProducts(params map[string]interface{}) ([]model.Product, error) {
	return s.productRepository.GetProducts(params)
}

func (s *productService) UpdateProduct(request *dto.ProductReq, productId int64) (dto.ProductSuccessResp, error) {
	_, err := s.productRepository.GetProductByID(productId)
	if err != nil {
		return dto.ProductSuccessResp{}, errs.NewErrDataNotFound("product id is not found", productId, errs.ErrorData{})
	}

	savedProduct, err := s.productRepository.UpdateProduct(request, productId)
	if err != nil {
		return dto.ProductSuccessResp{}, err
	}

	return dto.ProductSuccessResp{
		ID:        helper.IntToString(savedProduct.ID),
		CreatedAt: savedProduct.CreatedAtFormatter,
	}, nil
}

func (s *productService) DeleteProduct(productId int64) error {
	_, err := s.productRepository.GetProductByID(productId)
	if err != nil {
		return errs.NewErrDataNotFound("product id is not found", productId, errs.ErrorData{})
	}

	err = s.productRepository.DeleteProduct(productId)
	if err != nil {
		return err
	}

	return nil
}
