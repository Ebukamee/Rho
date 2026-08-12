package checkout

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rho-commerce/rho/internal/cart"
	"github.com/rho-commerce/rho/internal/discount"
	"github.com/rho-commerce/rho/internal/inventory"
	"github.com/rho-commerce/rho/internal/order"
	"github.com/rho-commerce/rho/internal/product"
)

var (
	ErrCartEmpty       = errors.New("cart is empty")
	ErrInvalidQuantity = errors.New("invalid cart quantity")
)

type Service struct {
	db          *pgxpool.Pool
	cartRepo    *cart.Repository
	productRepo *product.Repository
	discount    *discount.Service
}

func NewService(
	db *pgxpool.Pool,
	cartRepo *cart.Repository,
	productRepo *product.Repository,
	discountService *discount.Service,
) *Service {
	return &Service{
		db:          db,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		discount:    discountService,
	}
}

func (s *Service) Create(
	ctx context.Context,
	userID string,
	req Request,
) (*Response, error) {

	cartData, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart: %w", err)
	}

	if len(cartData.Items) == 0 {
		return nil, ErrCartEmpty
	}

	var (
		subtotal int64
		currency string
		items    []order.OrderItem
	)

	items = make([]order.OrderItem, 0, len(cartData.Items))

	for _, cartItem := range cartData.Items {
		if cartItem.Quantity < 1 {
			return nil, ErrInvalidQuantity
		}

		p, err := s.productRepo.GetByID(
			ctx,
			cartItem.ProductID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"get product %s: %w",
				cartItem.ProductID,
				err,
			)
		}

		if !p.Active {
			return nil, fmt.Errorf(
				"product %s is inactive",
				p.Name,
			)
		}

		if currency == "" {
			currency = p.Currency
		}

		if p.Currency != currency {
			return nil, errors.New(
				"cart contains products with different currencies",
			)
		}

		totalPrice := p.Price * int64(cartItem.Quantity)
		subtotal += totalPrice

		items = append(items, order.OrderItem{
			ID:         uuid.NewString(),
			ProductID:  p.ID,
			Name:       p.Name,
			SKU:        p.SKU,
			Quantity:   cartItem.Quantity,
			UnitPrice:  p.Price,
			TotalPrice: totalPrice,
		})
	}

	discountAmount := int64(0)

	if req.DiscountCode != "" {
		result, err := s.discount.Apply(
			ctx,
			discount.ApplyDiscountRequest{
				Code:  req.DiscountCode,
				Total: subtotal,
			},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"apply discount: %w",
				err,
			)
		}

		discountAmount = result.Discount
	}

	total := subtotal - discountAmount

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"begin checkout transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	txCartRepo := cart.NewRepositoryWithTx(tx)
	txInventoryRepo := inventory.NewRepositoryWithTx(tx)
	txOrderRepo := order.NewRepositoryWithTx(tx)

	// Re-read the cart inside the transaction so the transaction
	// works with the current cart state.
	currentCart, err := txCartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"get cart in transaction: %w",
			err,
		)
	}

	if len(currentCart.Items) == 0 {
		return nil, ErrCartEmpty
	}

	for _, item := range currentCart.Items {
		if err := txInventoryRepo.Reserve(
			ctx,
			item.ProductID,
			item.Quantity,
		); err != nil {
			return nil, fmt.Errorf(
				"reserve inventory: %w",
				err,
			)
		}
	}

	now := time.Now()

	newOrder := &order.Order{
		ID:        uuid.NewString(),
		UserID:    userID,
		Status:    order.OrderPending,
		Subtotal:  subtotal,
		Discount:  discountAmount,
		Total:     total,
		Currency:  currency,
		Items:     items,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := txOrderRepo.Create(
		ctx,
		newOrder,
	); err != nil {
		return nil, fmt.Errorf(
			"create order: %w",
			err,
		)
	}

	if err := txCartRepo.Clear(
		ctx,
		currentCart.ID,
	); err != nil {
		return nil, fmt.Errorf(
			"clear cart: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"commit checkout: %w",
			err,
		)
	}

	return &Response{
		OrderID:  newOrder.ID,
		Subtotal: subtotal,
		Discount: discountAmount,
		Total:    total,
		Currency: currency,
	}, nil
}

func (s *Service) Preview(
	ctx context.Context,
	userID string,
	req Request,
) (*Response, error) {
	cartData, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cartData.Items) == 0 {
		return nil, ErrCartEmpty
	}

	var subtotal int64
	var currency string

	for _, item := range cartData.Items {
		p, err := s.productRepo.GetByID(
			ctx,
			item.ProductID,
		)
		if err != nil {
			return nil, err
		}

		if !p.Active {
			return nil, fmt.Errorf(
				"product %s is inactive",
				p.Name,
			)
		}

		if currency == "" {
			currency = p.Currency
		}

		if p.Currency != currency {
			return nil, errors.New(
				"cart contains products with different currencies",
			)
		}

		subtotal += p.Price * int64(item.Quantity)
	}

	discountAmount := int64(0)

	if req.DiscountCode != "" {
		result, err := s.discount.Apply(
			ctx,
			discount.ApplyDiscountRequest{
				Code:  req.DiscountCode,
				Total: subtotal,
			},
		)
		if err != nil {
			return nil, err
		}

		discountAmount = result.Discount
	}

	return &Response{
		Subtotal: subtotal,
		Discount: discountAmount,
		Total:    subtotal - discountAmount,
		Currency: currency,
	}, nil
}
