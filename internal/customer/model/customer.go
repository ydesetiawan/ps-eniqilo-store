package model

type Customer struct {
	ID          int64  `json:"userId" db:"id"`
	Name        string `json:"name" db:"name"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
}
