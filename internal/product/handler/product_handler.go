package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/service"
	"ps-eniqilo-store/pkg/base/app"
	"ps-eniqilo-store/pkg/helper"
	"ps-eniqilo-store/pkg/httphelper/response"
	"strconv"
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

func (h *ProductHandler) UpdateProduct(ctx *app.Context) *response.WebResponse {
	vars := mux.Vars(ctx.Request)
	id, _ := vars["id"]
	productId, _ := strconv.Atoi(id)

	var request dto.ProductReq
	jsonString, _ := json.Marshal(ctx.GetJsonBody())
	err := json.Unmarshal(jsonString, &request)
	helper.Panic400IfError(err)

	err = dto.ValidateProductReq(request)
	helper.Panic400IfError(err)

	result, err := h.productService.UpdateProduct(&request, int64(productId))
	if err != nil {
		return &response.WebResponse{
			Status:  500,
			Message: err.Error(),
			Data:    result,
		}
	}

	return &response.WebResponse{
		Status:  200,
		Message: "successfully edit product",
		Data:    result,
	}
}

func (h *ProductHandler) DeleteProduct(ctx *app.Context) *response.WebResponse {
	vars := mux.Vars(ctx.Request)
	id, _ := vars["id"]
	productId, _ := strconv.Atoi(id)

	err := h.productService.DeleteProduct(int64(productId))
	if err != nil {
		return &response.WebResponse{
			Status:  500,
			Message: err.Error(),
		}
	}

	return &response.WebResponse{
		Status:  200,
		Message: "successfully delete product",
	}
}
