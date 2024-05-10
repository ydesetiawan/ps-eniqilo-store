package handler

import (
	"encoding/json"
	"ps-eniqilo-store/internal/customer/dto"
	"ps-eniqilo-store/internal/customer/model"
	"ps-eniqilo-store/internal/customer/service"
	"ps-eniqilo-store/pkg/base/app"
	"ps-eniqilo-store/pkg/helper"
	"ps-eniqilo-store/pkg/httphelper/response"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (h *CustomerHandler) CreateCustomer(ctx *app.Context) *response.WebResponse {
	var request dto.CustomerReq
	jsonString, _ := json.Marshal(ctx.GetJsonBody())
	err := json.Unmarshal(jsonString, &request)
	helper.Panic400IfError(err)

	err = dto.ValidateCustomerReq(request)
	helper.Panic400IfError(err)

	result, err := h.customerService.CreateCustomer(&request)

	helper.PanicIfError(err, "register user failed")

	return &response.WebResponse{
		Status:  201,
		Message: "successfully add customer",
		Data:    result,
	}
}

func (h *CustomerHandler) GetCustomers(ctx *app.Context) *response.WebResponse {
	reqParams, err := dto.GenerateCustomerReqParams(ctx)

	if err != nil {
		return &response.WebResponse{
			Status:  200,
			Message: err.Error(),
			Data:    []model.Customer{},
		}
	}

	result, err := h.customerService.SearchCustomers(reqParams)
	if err != nil {
		return &response.WebResponse{
			Status:  500,
			Message: err.Error(),
			Data:    result,
		}
	}

	return &response.WebResponse{
		Status: 200,
		Data:   result,
	}
}
