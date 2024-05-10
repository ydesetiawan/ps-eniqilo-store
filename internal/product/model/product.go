package model

import (
	"database/sql"
	"time"
)

type Product struct {
	ID                 int64        `json:"-" db:"id"`
	IDString           string       `json:"id" db:"-"`
	Name               string       `json:"name" db:"name"`
	SKU                string       `json:"sku" db:"sku"`
	Category           string       `json:"category" db:"category"`
	ImageURL           string       `json:"imageUrl" db:"image_url"`
	Notes              string       `json:"notes" db:"notes"`
	Price              float64      `json:"price" db:"price"`
	Stock              int          `json:"stock" db:"stock"`
	Location           string       `json:"location" db:"location"`
	IsAvailable        bool         `json:"isAvailable" db:"is_available"`
	CreatedAt          time.Time    `json:"-" db:"created_at"`
	CreatedAtFormatter string       `json:"createdAt"`
	DeletedAt          sql.NullTime `json:"-" db:"deleted_at"`
}
