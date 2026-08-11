package cart

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCart(ctx context.Context, userID string) (*Cart, error) {
	cart, err := s.repo.GetByUserID(ctx, userID)

	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			return s.createCart(ctx, userID)
		}

		return nil, fmt.Errorf("get cart: %w", err)
	}

	return cart, nil
}

func (s *Service) createCart(ctx context.Context, userID string) (*Cart, error) {
	now := time.Now()

	cart := &Cart{
		ID:        uuid.NewString(),
		UserID:    userID,
		Items:     []CartItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, cart); err != nil {
		if errors.Is(err, ErrDuplicateItem) {
			return s.repo.GetByUserID(ctx, userID)
		}

		return nil, fmt.Errorf("create cart: %w", err)
	}

	return cart, nil
}

func (s *Service) AddItem(
	ctx context.Context,
	userID string,
	req AddItemRequest,
) (*Cart, error) {
	productID := strings.TrimSpace(req.ProductID)

	if productID == "" {
		return nil, errors.New("product_id is required")
	}

	if req.Quantity < 1 {
		return nil, errors.New("quantity must be at least 1")
	}

	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	item := &CartItem{
		ID:        uuid.NewString(),
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  req.Quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.AddItem(ctx, item); err != nil {
		if errors.Is(err, ErrDuplicateItem) {
			return nil, ErrDuplicateItem
		}

		return nil, fmt.Errorf("add cart item: %w", err)
	}

	return s.GetCart(ctx, userID)
}

func (s *Service) UpdateItem(
	ctx context.Context,
	userID string,
	itemID string,
	req UpdateItemRequest,
) (*Cart, error) {
	if req.Quantity < 1 {
		return nil, errors.New("quantity must be at least 1")
	}

	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateItem(
		ctx,
		cart.ID,
		itemID,
		req.Quantity,
	); err != nil {
		return nil, fmt.Errorf("update cart item: %w", err)
	}

	return s.GetCart(ctx, userID)
}

func (s *Service) RemoveItem(
	ctx context.Context,
	userID string,
	itemID string,
) (*Cart, error) {
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.RemoveItem(ctx, cart.ID, itemID); err != nil {
		return nil, fmt.Errorf("remove cart item: %w", err)
	}

	return s.GetCart(ctx, userID)
}

func (s *Service) Clear(
	ctx context.Context,
	userID string,
) error {
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	if err := s.repo.Clear(ctx, cart.ID); err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}

	return nil
}
