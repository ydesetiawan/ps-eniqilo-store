package model

import (
	"ps-eniqilo-store/internal/user/dto"
	"time"
)

type User struct {
	ID          int64     `db:"id" json:"id"`
	PhoneNumber string    `db:"phone_number" json:"phoneNumber"`
	Password    string    `db:"password" json:"password"`
	Name        string    `db:"name" json:"name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func NewUser(req dto.RegisterReq) *User {
	return &User{
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		Name:        req.Name,
	}
}
