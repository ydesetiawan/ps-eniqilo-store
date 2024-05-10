package dto

import (
	"ps-eniqilo-store/pkg/base/app"

	"github.com/go-playground/validator/v10"
)

type CustomerReq struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=1,max=30,unique_phone_number"`
}

func ValidateCustomerReq(req CustomerReq) error {
	validate := validator.New()
	validate.RegisterValidation("unique_phone_number", func(fl validator.FieldLevel) bool {
		return true
	})

	return validate.Struct(req)
}

type CustomerResp struct {
	ID          string `json:"userId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

func GenerateCustomerReqParams(ctx *app.Context) (map[string]interface{}, error) {
	params := make(map[string]interface{})

	reqName := ctx.Request.URL.Query().Get("name")
	if reqName != "" {
		params["name"] = reqName
	}

	reqPhoneNumber := ctx.Request.URL.Query().Get("phoneNumber")
	if reqPhoneNumber != "" {
		params["phoneNumber"] = reqPhoneNumber
	}

	return params, nil
}
