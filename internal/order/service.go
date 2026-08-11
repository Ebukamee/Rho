package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	req CreateOrderRequest,
) (*Order, error) {

	now := time.Now()

	order := &Order{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		Status:    OrderPending,
		Subtotal:  req.Subtotal,
		Discount:  req.Discount,
		Total:     req.Total,
		Currency:  req.Currency,
		Items:     req.Items,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for i := range order.Items {
		order.Items[i].ID = uuid.NewString()
		order.Items[i].OrderID = order.ID
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	return order, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateStatus(
	ctx context.Context,
	id string,
	status OrderStatus,
) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
