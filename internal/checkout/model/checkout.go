package model

import "time"

type Checkout struct {
	ID                 int64     `json:"id" db:"id"`
	CustomerID         int64     `json:"customerId" db:"customer_id"`
	TotalPrice         float64   `json:"totalPrice" db:"total_price"`
	Paid               float64   `json:"paid" db:"paid"`
	Change             float64   `json:"change" db:"change"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	CreatedAtFormatter string    `json:"-"`
}
