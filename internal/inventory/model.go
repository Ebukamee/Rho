package inventory

import "time"

type Inventory struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Reserved  int       `json:"reserved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateInventoryRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"min=0"`
}

type UpdateInventoryRequest struct {
	Quantity int `json:"quantity" binding:"min=0"`
}

type AdjustInventoryRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

type InventoryResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Reserved  int       `json:"reserved"`
	Available int       `json:"available"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
