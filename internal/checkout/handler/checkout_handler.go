package handler

import (
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
