package dto

import (
	"ps-eniqilo-store/pkg/base/app"
	"ps-eniqilo-store/pkg/helper"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type ProductReq struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	SKU         string `json:"sku" validate:"required,min=1,max=30"`
	Category    string `json:"category" validate:"required,oneof=Clothing Accessories Footwear Beverages"`
	ImageURL    string `json:"imageUrl" validate:"required,imageURL"`
	Notes       string `json:"notes" validate:"required,min=1,max=201"`
	Price       int64  `json:"price" validate:"required,min=1"`
	Stock       int    `json:"stock" validate:"required,min=0,max=100000"`
	Location    string `json:"location" validate:"required,min=1,max=200"`
	IsAvailable *bool  `json:"isAvailable" validate:"boolean"`
}

func ValidateProductReq(req ProductReq) error {
	validate := validator.New()
	validate.RegisterValidation("imageURL", helper.ValidateURL)
	return validate.Struct(req)
}

type ProductSuccessResp struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

func isCategoryExists(val string) bool {
	categories := []string{
		"Clothing",
		"Accessories",
		"Footwear",
		"Beverages",
	}

	for _, c := range categories {
		if c == val {
			return true
		}
	}
	return false
}

func GenerateProductReqParams(ctx *app.Context) map[string]interface{} {
	params := make(map[string]interface{})

	reqProductId, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
	if err == nil {
		params["id"] = reqProductId
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

	reqName := ctx.Request.URL.Query().Get("name")
	if "" != reqName {
		params["name"] = reqName
	}

	reqCategory := ctx.Request.URL.Query().Get("category")
	if "" != reqCategory {
		if isCategoryExists(reqCategory) {
			params["category"] = reqCategory
		}
	}

	reqIsAvailable := ctx.Request.URL.Query().Get("isAvailable")
	isAvailable, err := strconv.ParseBool(reqIsAvailable)
	if err == nil {
		params["isAvailable"] = isAvailable
	}

	reqSku := ctx.Request.URL.Query().Get("sku")
	if "" != reqSku {
		params["sku"] = reqSku
	}

	reqPriceOrderBy := ctx.Request.URL.Query().Get("price")
	if "" != reqPriceOrderBy && helper.IsOrderValueValid(reqPriceOrderBy) {
		params["price"] = reqPriceOrderBy
	}

	reqInStock := ctx.Request.URL.Query().Get("inStock")
	inStock, err := strconv.ParseBool(reqInStock)
	if err == nil {
		params["inStock"] = inStock
	}

	reqCreatedAtOrderBy := ctx.Request.URL.Query().Get("createdAt")
	if "" != reqCreatedAtOrderBy && helper.IsOrderValueValid(reqCreatedAtOrderBy) {
		params["createdAt"] = reqCreatedAtOrderBy
	} else {
		params["createdAt"] = "desc"
	}

	return params
}

func GenerateSearchSKUReqParams(ctx *app.Context) map[string]interface{} {
	params := make(map[string]interface{})

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

	reqName := ctx.Request.URL.Query().Get("name")
	if "" != reqName {
		params["name"] = reqName
	}

	reqCategory := ctx.Request.URL.Query().Get("category")
	if "" != reqCategory {
		if isCategoryExists(reqCategory) {
			params["category"] = reqCategory
		}
	}

	params["isAvailable"] = true

	reqSku := ctx.Request.URL.Query().Get("sku")
	if "" != reqSku {
		params["sku"] = reqSku
	}

	reqPriceOrderBy := ctx.Request.URL.Query().Get("price")
	if "" != reqPriceOrderBy && helper.IsOrderValueValid(reqPriceOrderBy) {
		params["price"] = reqPriceOrderBy
	}

	reqInStock := ctx.Request.URL.Query().Get("inStock")
	inStock, err := strconv.ParseBool(reqInStock)
	if err == nil {
		params["inStock"] = inStock
	}

	params["createdAt"] = "desc"

	return params
}
