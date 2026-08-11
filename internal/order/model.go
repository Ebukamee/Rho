package order

import "time"

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderPaid       OrderStatus = "paid"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Status    OrderStatus `json:"status"`
	Subtotal  int64       `json:"subtotal"`
	Discount  int64       `json:"discount"`
	Total     int64       `json:"total"`
	Currency  string      `json:"currency"`
	Items     []OrderItem `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID         string `json:"id"`
	OrderID    string `json:"order_id"`
	ProductID  string `json:"product_id"`
	Name       string `json:"name"`
	SKU        string `json:"sku"`
	Quantity   int    `json:"quantity"`
	UnitPrice  int64  `json:"unit_price"`
	TotalPrice int64  `json:"total_price"`
}

type CreateOrderRequest struct {
	UserID   string      `json:"user_id"`
	Subtotal int64       `json:"subtotal"`
	Discount int64       `json:"discount"`
	Total    int64       `json:"total"`
	Currency string      `json:"currency"`
	Items    []OrderItem `json:"items"`
}

type UpdateStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required"`
}
