package user

import "time"

type Address struct {
	ID          string    `json:"id"`
	UserID      string    `json:"-"`
	Label       string    `json:"label"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Company     string    `json:"company,omitempty"`
	Phone       string    `json:"phone"`
	Line1       string    `json:"line1"`
	Line2       string    `json:"line2,omitempty"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	PostalCode  string    `json:"postal_code"`
	CountryCode string    `json:"country_code"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateAddressRequest struct {
	Label       string `json:"label" binding:"omitempty,max=50"`
	FirstName   string `json:"first_name" binding:"required,max=100"`
	LastName    string `json:"last_name" binding:"required,max=100"`
	Company     string `json:"company" binding:"omitempty,max=100"`
	Phone       string `json:"phone" binding:"required,max=30"`
	Line1       string `json:"line1" binding:"required,max=255"`
	Line2       string `json:"line2" binding:"omitempty,max=255"`
	City        string `json:"city" binding:"required,max=100"`
	State       string `json:"state" binding:"required,max=100"`
	PostalCode  string `json:"postal_code" binding:"required,max=20"`
	CountryCode string `json:"country_code" binding:"required,len=2"`
}

type UpdateAddressRequest struct {
	Label       *string `json:"label" binding:"omitempty,max=50"`
	FirstName   *string `json:"first_name" binding:"omitempty,max=100"`
	LastName    *string `json:"last_name" binding:"omitempty,max=100"`
	Company     *string `json:"company" binding:"omitempty,max=100"`
	Phone       *string `json:"phone" binding:"omitempty,max=30"`
	Line1       *string `json:"line1" binding:"omitempty,max=255"`
	Line2       *string `json:"line2" binding:"omitempty,max=255"`
	City        *string `json:"city" binding:"omitempty,max=100"`
	State       *string `json:"state" binding:"omitempty,max=100"`
	PostalCode  *string `json:"postal_code" binding:"omitempty,max=20"`
	CountryCode *string `json:"country_code" binding:"omitempty,len=2"`
	IsDefault   *bool   `json:"is_default"`
}
