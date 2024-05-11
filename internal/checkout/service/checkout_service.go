package service

import (
	"context"
	"ps-eniqilo-store/internal/checkout/dto"
	"ps-eniqilo-store/internal/checkout/model"
	"ps-eniqilo-store/internal/checkout/repository"
	customerrepo "ps-eniqilo-store/internal/customer/repository"
	productmodel "ps-eniqilo-store/internal/product/model"
	productrepo "ps-eniqilo-store/internal/product/repository"
	"ps-eniqilo-store/pkg/errs"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type CheckoutService interface {
	GetCheckOutHistory(map[string]interface{}) ([]dto.CheckOutHistoryResp, error)
	ProductCheckout(context.Context, *dto.ProductCheckoutReq) error
}

type checkoutService struct {
	db                 *sqlx.DB
	checkoutRepository repository.CheckoutRepository
	productRepository  productrepo.ProductRepository
	customerRepository customerrepo.CustomerRepository
}

func NewCheckoutServiceImpl(
	checkoutRepository repository.CheckoutRepository,
	productRepository productrepo.ProductRepository,
	customerRepository customerrepo.CustomerRepository,
	db *sqlx.DB,
) CheckoutService {
	return &checkoutService{
		checkoutRepository: checkoutRepository,
		productRepository:  productRepository,
		customerRepository: customerRepository,
		db:                 db,
	}
}

func (c *checkoutService) GetCheckOutHistory(request map[string]interface{}) ([]dto.CheckOutHistoryResp, error) {
	return c.checkoutRepository.GetCheckoutHistory(request)
}

/*
- `200` successfully checkout product
- `404` `customerId` is not found
- `404` one of productIds is not found v
- `400` request doesn’t pass validation
- `400` `paid` is not enough based on all bought product v
- `400` `change` is not right, based on all bought product, and what is paid v
- `400` one of productIds stock is not enough
- `400` one of productIds `isAvailable == false`
*/
func (c *checkoutService) ProductCheckout(ctx context.Context, request *dto.ProductCheckoutReq) (err error) {
	productDetails, err := c.getProductDetails(request)
	if err != nil {
		return
	}

	err = c.validateCheckout(request, productDetails)
	if err != nil {
		return
	}

	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	customerId, _ := strconv.ParseInt(request.CustomerId, 10, 64)
	checkout := &model.Checkout{
		CustomerID: customerId,
		TotalPrice: request.Paid - *request.Change,
		Paid:       request.Paid,
		Change:     *request.Change,
	}

	checkoutId, err := c.checkoutRepository.CreateCheckout(tx, checkout)
	if err != nil {
		return err
	}

	for _, pd := range request.ProductDetails {
		productId, _ := strconv.ParseInt(pd.ProductID, 10, 64)
		product := productDetails[productId]

		checkoutDetail := &model.CheckoutDetail{
			CheckoutID:   checkoutId,
			ProductID:    productId,
			ProductPrice: product.Price,
			TotalPrice:   product.Price * float64(pd.Quantity),
			Quantity:     pd.Quantity,
		}
		err := c.checkoutRepository.CreateCheckoutDetail(tx, checkoutDetail)
		if err != nil {
			return err
		}

		err = c.productRepository.DecreaseStock(tx, productId, pd.Quantity)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (c *checkoutService) getProductDetails(request *dto.ProductCheckoutReq) (map[int64]productmodel.Product, error) {
	productIds, err := getUniqueProductIds(request.ProductDetails)
	if err != nil {
		return nil, err
	}

	products, err := c.productRepository.GetProductByIDs(productIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[int64]productmodel.Product)
	for _, pd := range products {
		productMap[pd.ID] = pd
	}

	return productMap, nil
}

func (c *checkoutService) validateCheckout(request *dto.ProductCheckoutReq, productDetails map[int64]productmodel.Product) (err error) {
	if request.Change == nil {
		return errs.NewErrBadRequest("empty change")
	}

	customerId, err := strconv.ParseInt(request.CustomerId, 10, 64)
	if err != nil {
		return errs.NewErrDataNotFound("invalid customer id", customerId, errs.ErrorData{})
	}

	_, err = c.customerRepository.GetCustomerByID(customerId)
	if err != nil {
		return errs.NewErrDataNotFound("invalid customer id", customerId, errs.ErrorData{})
	}

	totalPrice := float64(0)
	for _, pd := range request.ProductDetails {
		productId, err := strconv.ParseInt(pd.ProductID, 10, 64)
		if err != nil {
			return errs.NewErrDataNotFound("invalid product id", productId, errs.ErrorData{})
		}

		product, ok := productDetails[productId]
		if !ok {
			return errs.NewErrDataNotFound("product is not found", productId, errs.ErrorData{})
		}
		totalPrice += product.Price * float64(pd.Quantity)

		if !product.IsAvailable {
			return errs.NewErrBadRequest("product is not available")
		}

		if product.Stock < pd.Quantity {
			return errs.NewErrBadRequest("product stock is not enough")
		}
	}

	if request.Paid < totalPrice {
		return errs.NewErrBadRequest("paid price is not enough")
	}

	if (request.Paid - *request.Change) != totalPrice {
		return errs.NewErrBadRequest("change is not valid based on paid price")
	}

	return
}

func getUniqueProductIds(productDetails []dto.ProductDetail) ([]int64, error) {
	uniqueIDs := make(map[string]bool)
	for _, pd := range productDetails {
		uniqueIDs[pd.ProductID] = true
	}

	ids := make([]int64, 0, len(uniqueIDs))
	for idStr := range uniqueIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, errs.NewErrDataNotFound("invalid product id", idStr, errs.ErrorData{})
		}
		ids = append(ids, id)
	}
	return ids, nil
}
