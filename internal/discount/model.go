package discount

import "time"

type DiscountType string

const (
	Percentage DiscountType = "percentage"
	Fixed      DiscountType = "fixed"
)

type Discount struct {
	ID           string       `json:"id"`
	Code         string       `json:"code"`
	Type         DiscountType `json:"type"`
	Value        int64        `json:"value"`
	MinimumOrder int64        `json:"minimum_order"`
	UsageLimit   *int         `json:"usage_limit,omitempty"`
	UsageCount   int          `json:"usage_count"`
	StartsAt     *time.Time   `json:"starts_at,omitempty"`
	ExpiresAt    *time.Time   `json:"expires_at,omitempty"`
	Active       bool         `json:"active"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type CreateDiscountRequest struct {
	Code         string       `json:"code" binding:"required"`
	Type         DiscountType `json:"type" binding:"required"`
	Value        int64        `json:"value" binding:"required,min=1"`
	MinimumOrder int64        `json:"minimum_order" binding:"min=0"`
	UsageLimit   *int         `json:"usage_limit"`
	StartsAt     *time.Time   `json:"starts_at"`
	ExpiresAt    *time.Time   `json:"expires_at"`
	Active       *bool        `json:"active"`
}

type UpdateDiscountRequest struct {
	Code         string       `json:"code"`
	Type         DiscountType `json:"type"`
	Value        *int64       `json:"value"`
	MinimumOrder *int64       `json:"minimum_order"`
	UsageLimit   *int         `json:"usage_limit"`
	StartsAt     *time.Time   `json:"starts_at"`
	ExpiresAt    *time.Time   `json:"expires_at"`
	Active       *bool        `json:"active"`
}

type ApplyDiscountRequest struct {
	Code  string `json:"code" binding:"required"`
	Total int64  `json:"total" binding:"min=0"`
}

type DiscountResult struct {
	DiscountID    string `json:"discount_id"`
	Code          string `json:"code"`
	OriginalTotal int64  `json:"original_total"`
	Discount      int64  `json:"discount"`
	FinalTotal    int64  `json:"final_total"`
}
