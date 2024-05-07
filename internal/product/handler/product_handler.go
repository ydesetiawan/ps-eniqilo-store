package handler

import (
	"ps-eniqilo-store/internal/product/service"
	"ps-eniqilo-store/pkg/base/app"
	"ps-eniqilo-store/pkg/httphelper/response"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(ctx *app.Context) *response.WebResponse {
	return nil
}
