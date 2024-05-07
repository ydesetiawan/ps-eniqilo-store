package handler

import (
	"encoding/json"
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/service"
	"ps-eniqilo-store/pkg/base/app"
	"ps-eniqilo-store/pkg/helper"
	"ps-eniqilo-store/pkg/httphelper/response"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(ctx *app.Context) *response.WebResponse {
	var request dto.ProductReq
	jsonString, _ := json.Marshal(ctx.GetJsonBody())
	err := json.Unmarshal(jsonString, &request)
	helper.Panic400IfError(err)

	err = dto.ValidateProductReq(request)
	helper.Panic400IfError(err)

	result, err := h.productService.CreateProduct(&request)
	if err != nil {
		return &response.WebResponse{
			Status:  500,
			Message: err.Error(),
			Data:    result,
		}
	}

	return &response.WebResponse{
		Status:  201,
		Message: "successfully add product",
		Data:    result,
	}
}
