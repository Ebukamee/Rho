package inventory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("quantity cannot be negative")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	req CreateInventoryRequest,
) (*Inventory, error) {
	if req.Quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	now := time.Now()

	inventory := &Inventory{
		ID:        uuid.NewString(),
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Reserved:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, inventory); err != nil {
		return nil, fmt.Errorf("create inventory: %w", err)
	}

	return inventory, nil
}

func (s *Service) GetByProductID(
	ctx context.Context,
	productID string,
) (*InventoryResponse, error) {
	inventory, err := s.repo.GetByProductID(ctx, productID)

	if err != nil {
		return nil, err
	}

	return toResponse(inventory), nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (*InventoryResponse, error) {
	inventory, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return toResponse(inventory), nil
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	req UpdateInventoryRequest,
) (*InventoryResponse, error) {
	if req.Quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	inventory, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if req.Quantity < inventory.Reserved {
		return nil, errors.New(
			"quantity cannot be less than reserved stock",
		)
	}

	inventory.Quantity = req.Quantity
	inventory.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("update inventory: %w", err)
	}

	return toResponse(inventory), nil
}

func (s *Service) Adjust(
	ctx context.Context,
	productID string,
	quantity int,
) (*InventoryResponse, error) {
	inventory, err := s.repo.GetByProductID(ctx, productID)

	if err != nil {
		return nil, err
	}

	if inventory.Quantity+quantity < inventory.Reserved {
		return nil, ErrInvalidQuantity
	}

	if err := s.repo.Adjust(ctx, productID, quantity); err != nil {
		return nil, fmt.Errorf("adjust inventory: %w", err)
	}

	return s.GetByProductID(ctx, productID)
}

func (s *Service) Reserve(
	ctx context.Context,
	productID string,
	quantity int,
) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	err := s.repo.Reserve(ctx, productID, quantity)

	if errors.Is(err, ErrNotFound) {
		return ErrInsufficientStock
	}

	if err != nil {
		return fmt.Errorf("reserve inventory: %w", err)
	}

	return nil
}

func (s *Service) Release(
	ctx context.Context,
	productID string,
	quantity int,
) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if err := s.repo.Release(ctx, productID, quantity); err != nil {
		return fmt.Errorf("release inventory: %w", err)
	}

	return nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete inventory: %w", err)
	}

	return nil
}

func toResponse(inventory *Inventory) *InventoryResponse {
	return &InventoryResponse{
		ID:        inventory.ID,
		ProductID: inventory.ProductID,
		Quantity:  inventory.Quantity,
		Reserved:  inventory.Reserved,
		Available: inventory.Quantity - inventory.Reserved,
		CreatedAt: inventory.CreatedAt,
		UpdatedAt: inventory.UpdatedAt,
	}
}
