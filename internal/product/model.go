package product

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	SKU         string    `json:"sku"`
	Price       int64     `json:"price"`
	Currency    string    `json:"currency"`
	ImageURL    string    `json:"image_url"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	SKU         string `json:"sku" binding:"required"`
	Price       int64  `json:"price" binding:"gte=0"`
	Currency    string `json:"currency" binding:"required,len=3"`
	ImageURL    string `json:"image_url"`
	Active      *bool  `json:"active"`
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	SKU         string `json:"sku"`
	Price       *int64 `json:"price"`
	Currency    string `json:"currency"`
	ImageURL    string `json:"image_url"`
	Active      *bool  `json:"active"`
}

type ProductListResponse struct {
	Products   []Product `json:"products"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}
