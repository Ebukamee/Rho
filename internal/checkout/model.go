package checkout

type Request struct {
	DiscountCode string `json:"discount_code"`
}

type Response struct {
	OrderID  string `json:"order_id"`
	Subtotal int64  `json:"subtotal"`
	Discount int64  `json:"discount"`
	Total    int64  `json:"total"`
	Currency string `json:"currency"`
}
