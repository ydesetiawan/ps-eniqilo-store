package model

import (
	"database/sql"
	"time"
)

type Product struct {
	ID                 int64        `json:"id" db:"id"`
	Name               string       `json:"name" db:"name"`
	SKU                string       `json:"sku" db:"sku"`
	Category           string       `json:"category" db:"category"`
	ImageURL           string       `json:"imageUrl" db:"image_url"`
	Notes              string       `json:"notes" db:"notes"`
	Price              int64        `json:"price" db:"price"`
	Stock              int          `json:"stock" db:"stock"`
	Location           string       `json:"location" db:"location"`
	IsAvailable        bool         `json:"isAvailable" db:"is_available"`
	CreatedAt          time.Time    `json:"createdAt" db:"created_at"`
	CreatedAtFormatter string       `json:"-"`
	DeletedAt          sql.NullTime `json:"-" db:"deleted_at"`
}
