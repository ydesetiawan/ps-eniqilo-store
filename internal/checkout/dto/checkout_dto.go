package dto

import (
	"ps-eniqilo-store/pkg/base/app"
	"ps-eniqilo-store/pkg/helper"
	"strconv"
)

type CheckOutHistoryResp struct {
	TransactionId  string          `json:"transaction_id"`
	CustomerId     string          `json:"customer_id"`
	ProductDetails []ProductDetail `json:"productDetails"`
	Paid           float64         `json:"paid"`
	Change         float64         `json:"change"`
	CreatedAt      string          `json:"createdAt"`
}

type ProductDetail struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func GenerateCheckoutHistoryReqParams(ctx *app.Context) map[string]interface{} {
	params := make(map[string]interface{})

	reqProductId, err := strconv.Atoi(ctx.Request.URL.Query().Get("customerId"))
	if err == nil {
		params["customerId"] = reqProductId
	}

	reqLimit, err := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
	if err != nil {
		reqLimit = 5
	}
	params["limit"] = reqLimit

	reqOffset, err := strconv.Atoi(ctx.Request.URL.Query().Get("offset"))
	if err != nil {
		reqOffset = 0
	}
	params["offset"] = reqOffset

	reqCreatedAtOrderBy := ctx.Request.URL.Query().Get("createdAt")
	if "" != reqCreatedAtOrderBy && helper.IsOrderValueValid(reqCreatedAtOrderBy) {
		params["createdAt"] = reqCreatedAtOrderBy
	}

	return params
}
