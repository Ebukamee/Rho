package shipping

import "time"

type ShipmentStatus string

const (
	StatusPending   ShipmentStatus = "pending"
	StatusProcessing ShipmentStatus = "processing"
	StatusShipped   ShipmentStatus = "shipped"
	StatusDelivered ShipmentStatus = "delivered"
	StatusCancelled ShipmentStatus = "cancelled"
)

type Shipment struct {
	ID            string         `json:"id"`
	OrderID       string         `json:"order_id"`
	UserID        string         `json:"user_id"`
	Carrier       string         `json:"carrier"`
	Service      string         `json:"service"`
	TrackingNumber string        `json:"tracking_number,omitempty"`
	Cost          int64          `json:"cost"`
	Currency      string         `json:"currency"`
	Status        ShipmentStatus `json:"status"`
	EstimatedDays int            `json:"estimated_days"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type CreateShipmentRequest struct {
	OrderID        string `json:"order_id" binding:"required"`
	Carrier        string `json:"carrier" binding:"required"`
	Service        string `json:"service" binding:"required"`
	TrackingNumber string `json:"tracking_number"`
	Cost           int64  `json:"cost" binding:"gte=0"`
	Currency       string `json:"currency" binding:"required"`
	EstimatedDays  int    `json:"estimated_days" binding:"gte=0"`
}

type UpdateShipmentRequest struct {
	Carrier        string         `json:"carrier"`
	Service        string         `json:"service"`
	TrackingNumber string         `json:"tracking_number"`
	Status         ShipmentStatus `json:"status"`
	EstimatedDays  *int           `json:"estimated_days"`
}