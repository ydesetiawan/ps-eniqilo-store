package model

import "time"

type CheckoutDetail struct {
	ID                 int64     `json:"id" db:"id"`
	CheckoutID         int64     `json:"checkoutId" db:"checkout_id"`
	ProductID          int64     `json:"productId" db:"product_id"`
	ProductPrice       float64   `json:"productPrice" db:"product_price"`
	TotalPrice         float64   `json:"totalPrice" db:"total_price"`
	Quantity           int       `json:"quantity" db:"quantity"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	CreatedAtFormatter string    `json:"-"`
}
