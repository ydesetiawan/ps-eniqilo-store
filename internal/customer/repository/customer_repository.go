package repository

import (
	"ps-eniqilo-store/internal/customer/dto"
	"ps-eniqilo-store/internal/customer/model"
)

type CustomerRepository interface {
	CustomerEmailExists(string) (bool, error)
	GetCustomerByID(int64) (model.Customer, error)
	CreateCustomer(*dto.CustomerReq) (model.Customer, error)
	SearchCustomers(params map[string]interface{}) ([]model.Customer, error)
}
