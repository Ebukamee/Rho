package payment

import "time"

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID          string        `json:"id"`
	OrderID     string        `json:"order_id"`
	UserID      string        `json:"user_id"`
	Provider    string        `json:"provider"`
	ProviderRef string        `json:"provider_ref,omitempty"`
	Amount      int64         `json:"amount"`
	Currency    string        `json:"currency"`
	Status      PaymentStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type InitializePaymentRequest struct {
	OrderID  string `json:"order_id" binding:"required"`
	Provider string `json:"provider" binding:"required"`
}

type PaymentInitialization struct {
	PaymentID        string `json:"payment_id"`
	Provider         string `json:"provider"`
	ProviderRef      string `json:"provider_ref"`
	AuthorizationURL string `json:"authorization_url"`
}

type PaymentVerification struct {
	ProviderRef string        `json:"provider_ref"`
	Status      PaymentStatus `json:"status"`
}
