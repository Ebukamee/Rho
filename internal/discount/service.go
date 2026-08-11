package discount

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidType       = errors.New("invalid discount type")
	ErrInvalidValue      = errors.New("invalid discount value")
	ErrInactive          = errors.New("discount is inactive")
	ErrExpired           = errors.New("discount has expired")
	ErrNotStarted        = errors.New("discount is not active yet")
	ErrMinimumNotReached = errors.New("minimum order amount not reached")
	ErrUsageLimit        = errors.New("discount usage limit reached")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	req CreateDiscountRequest,
) (*Discount, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))

	if code == "" {
		return nil, errors.New("code is required")
	}

	if req.Type != Percentage && req.Type != Fixed {
		return nil, ErrInvalidType
	}

	if req.Value <= 0 {
		return nil, ErrInvalidValue
	}

	if req.Type == Percentage && req.Value > 100 {
		return nil, errors.New("percentage cannot exceed 100")
	}

	if req.UsageLimit != nil && *req.UsageLimit < 1 {
		return nil, errors.New("usage limit must be greater than zero")
	}

	now := time.Now()

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	d := &Discount{
		ID:           uuid.NewString(),
		Code:         code,
		Type:         req.Type,
		Value:        req.Value,
		MinimumOrder: req.MinimumOrder,
		UsageLimit:   req.UsageLimit,
		StartsAt:     req.StartsAt,
		ExpiresAt:    req.ExpiresAt,
		Active:       active,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("create discount: %w", err)
	}

	return d, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (*Discount, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	req UpdateDiscountRequest,
) (*Discount, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Code != "" {
		d.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	}

	if req.Type != "" {
		if req.Type != Percentage && req.Type != Fixed {
			return nil, ErrInvalidType
		}

		d.Type = req.Type
	}

	if req.Value != nil {
		if *req.Value <= 0 {
			return nil, ErrInvalidValue
		}

		if d.Type == Percentage && *req.Value > 100 {
			return nil, ErrInvalidValue
		}

		d.Value = *req.Value
	}

	if req.MinimumOrder != nil {
		if *req.MinimumOrder < 0 {
			return nil, ErrInvalidValue
		}

		d.MinimumOrder = *req.MinimumOrder
	}

	if req.UsageLimit != nil {
		if *req.UsageLimit < 1 {
			return nil, ErrInvalidValue
		}

		d.UsageLimit = req.UsageLimit
	}

	if req.StartsAt != nil {
		d.StartsAt = req.StartsAt
	}

	if req.ExpiresAt != nil {
		d.ExpiresAt = req.ExpiresAt
	}

	if req.Active != nil {
		d.Active = *req.Active
	}

	d.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, d); err != nil {
		return nil, fmt.Errorf("update discount: %w", err)
	}

	return d, nil
}

func (s *Service) Apply(
	ctx context.Context,
	req ApplyDiscountRequest,
) (*DiscountResult, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))

	d, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if !d.Active {
		return nil, ErrInactive
	}

	if d.StartsAt != nil && now.Before(*d.StartsAt) {
		return nil, ErrNotStarted
	}

	if d.ExpiresAt != nil && now.After(*d.ExpiresAt) {
		return nil, ErrExpired
	}

	if req.Total < d.MinimumOrder {
		return nil, ErrMinimumNotReached
	}

	if d.UsageLimit != nil && d.UsageCount >= *d.UsageLimit {
		return nil, ErrUsageLimit
	}

	var amount int64

	switch d.Type {
	case Percentage:
		amount = req.Total * d.Value / 100

	case Fixed:
		amount = d.Value

	default:
		return nil, ErrInvalidType
	}

	if amount > req.Total {
		amount = req.Total
	}

	return &DiscountResult{
		DiscountID:    d.ID,
		Code:          d.Code,
		OriginalTotal: req.Total,
		Discount:      amount,
		FinalTotal:    req.Total - amount,
	}, nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete discount: %w", err)
	}

	return nil
}
