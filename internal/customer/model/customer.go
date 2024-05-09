package model

import "time"

type Customer struct {
	ID                 int64     `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	PhoneNumber        string    `json:"phoneNumber" db:"phone_number"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	CreatedAtFormatter string    `json:"-"`
}
