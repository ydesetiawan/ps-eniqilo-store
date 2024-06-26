package handler

import (
	"encoding/json"
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/model"
	"ps-eniqilo-store/internal/product/service"
	"ps-eniqilo-store/pkg/base/app"
	"ps-eniqilo-store/pkg/helper"
	"ps-eniqilo-store/pkg/httphelper/response"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetProduct(ctx *app.Context) *response.WebResponse {
	reqParams := dto.GenerateProductReqParams(ctx)

	results, err := h.productService.GetProducts(reqParams)
	helper.PanicIfError(err, "error when [GetProduct")

	message := "successfully get product"
	if len(results) == 0 {
		message = "DATA NOT FOUND"
		results = []model.Product{}
	}

	return &response.WebResponse{
		Status:  200,
		Message: message,
		Data:    results,
	}
}

func (h *ProductHandler) SearchSKU(ctx *app.Context) *response.WebResponse {
	reqParams := dto.GenerateSearchSKUReqParams(ctx)

	results, err := h.productService.GetProducts(reqParams)
	helper.PanicIfError(err, "error when [GetProduct")

	message := "successfully get product"
	if len(results) == 0 {
		message = "DATA NOT FOUND"
	}

	return &response.WebResponse{
		Status:  200,
		Message: message,
		Data:    results,
	}
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
	id := vars["id"]
	if !helper.IdIsInteger(id) {
		return &response.WebResponse{
			Status:  404,
			Message: "product id is invalid",
			Data:    nil,
		}
	}

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
	id := vars["id"]
	if !helper.IdIsInteger(id) {
		return &response.WebResponse{
			Status:  404,
			Message: "product id is invalid",
			Data:    nil,
		}
	}

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
