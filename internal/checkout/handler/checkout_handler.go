package handler

import (
	"encoding/json"
	"ps-eniqilo-store/internal/checkout/dto"
	"ps-eniqilo-store/internal/checkout/service"
	"ps-eniqilo-store/pkg/base/app"
	"ps-eniqilo-store/pkg/helper"
	"ps-eniqilo-store/pkg/httphelper/response"
)

type CheckoutHandler struct {
	checkoutService service.CheckoutService
}

func NewCheckoutHandler(checkoutService service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{checkoutService: checkoutService}
}

func (h *CheckoutHandler) CheckoutProduct(ctx *app.Context) *response.WebResponse {
	var request dto.ProductCheckoutReq
	jsonString, _ := json.Marshal(ctx.GetJsonBody())
	err := json.Unmarshal(jsonString, &request)
	helper.Panic400IfError(err)

	err = helper.ValidateStruct(request)
	helper.Panic400IfError(err)

	err = h.checkoutService.ProductCheckout(ctx.Context(), &request)
	helper.PanicIfError(err, "error CheckoutProduct")

	return &response.WebResponse{
		Status:  200,
		Message: "successfully checkout product",
	}
}

func (h *CheckoutHandler) GetCheckoutHistory(ctx *app.Context) *response.WebResponse {
	reqParams := dto.GenerateCheckoutHistoryReqParams(ctx)

	results, err := h.checkoutService.GetCheckOutHistory(reqParams)
	helper.PanicIfError(err, "error when [GetCheckoutHistory")

	message := "successfully get checkout history"
	if len(results) == 0 {
		message = "DATA NOT FOUND"
	}

	return &response.WebResponse{
		Status:  200,
		Message: message,
		Data:    results,
	}
}
