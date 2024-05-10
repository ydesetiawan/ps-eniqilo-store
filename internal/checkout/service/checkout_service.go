package service

import (
	"ps-eniqilo-store/internal/checkout/dto"
	"ps-eniqilo-store/internal/checkout/repository"
)

type CheckoutService interface {
	GetCheckOutHistory(map[string]interface{}) ([]dto.CheckOutHistoryResp, error)
}

type checkoutService struct {
	checkoutRepository repository.CheckoutRepository
}

func NewCheckoutServiceImpl(checkoutRepository repository.CheckoutRepository) CheckoutService {
	return &checkoutService{checkoutRepository: checkoutRepository}
}

func (c checkoutService) GetCheckOutHistory(request map[string]interface{}) ([]dto.CheckOutHistoryResp, error) {
	return c.checkoutRepository.GetCheckoutHistory(request)
}
