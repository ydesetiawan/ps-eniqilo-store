package service

import (
	"ps-eniqilo-store/internal/customer/dto"
	"ps-eniqilo-store/internal/customer/repository"
	"ps-eniqilo-store/pkg/errs"
	"ps-eniqilo-store/pkg/helper"
)

type CustomerService interface {
	GetCustomerByID(int64) (dto.CustomerResp, error)
	CreateCustomer(*dto.CustomerReq) (dto.CustomerResp, error)
	SearchCustomers(map[string]interface{}) ([]dto.CustomerResp, error)
}

type customerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerServiceImpl(customerRepository repository.CustomerRepository) CustomerService {
	return &customerService{customerRepository: customerRepository}
}

func (s *customerService) CreateCustomer(customer *dto.CustomerReq) (dto.CustomerResp, error) {
	savedCustomer, err := s.customerRepository.CreateCustomer(customer)

	if err != nil {
		return dto.CustomerResp{}, errs.NewErrDataConflict("phone number is already exist", customer.PhoneNumber)
	}

	return dto.CustomerResp{
		ID:          helper.IntToString(savedCustomer.ID),
		Name:        savedCustomer.Name,
		PhoneNumber: savedCustomer.PhoneNumber,
	}, nil
}

func (s *customerService) GetCustomerByID(customerId int64) (dto.CustomerResp, error) {
	_, err := s.customerRepository.GetCustomerByID(customerId)
	if err != nil {
		return dto.CustomerResp{}, errs.NewErrDataNotFound("customer id is not found", customerId, errs.ErrorData{})
	}

	savedCustomer, err := s.customerRepository.GetCustomerByID(customerId)
	if err != nil {
		return dto.CustomerResp{}, err
	}

	return dto.CustomerResp{
		ID:          helper.IntToString(savedCustomer.ID),
		Name:        savedCustomer.Name,
		PhoneNumber: savedCustomer.PhoneNumber,
	}, nil
}

func (s *customerService) SearchCustomers(params map[string]interface{}) ([]dto.CustomerResp, error) {
	customers, err := s.customerRepository.SearchCustomers(params)
	if err != nil {
		return []dto.CustomerResp{}, err
	}

	var customerResp []dto.CustomerResp
	for _, customer := range customers {
		customerResp = append(customerResp, dto.CustomerResp{
			ID:          helper.IntToString(customer.ID),
			Name:        customer.Name,
			PhoneNumber: customer.PhoneNumber,
		})
	}

	return customerResp, nil
}
