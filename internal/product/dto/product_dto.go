package dto

type ProductReq struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	SKU         string `json:"sku" validate:"required,min=1,max=30"`
	Category    string `json:"category" validate:"required,oneof=Clothing Accessories Footwear Beverages"`
	ImageURL    string `json:"imageUrl" validate:"required,url"`
	Notes       string `json:"notes" validate:"required,min=1,max=200"`
	Price       int64  `json:"price" validate:"required,min=1"`
	Stock       int    `json:"stock" validate:"required,min=0,max=100000"`
	Location    string `json:"location" validate:"required,min=1,max=200"`
	IsAvailable bool   `json:"isAvailable" validate:"required"`
}

type ProductRes struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}
