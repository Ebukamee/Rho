package payment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rho-commerce/rho/internal/order"
)

var (
	ErrProviderNotFound = errors.New(
		"payment provider not found",
	)

	ErrInvalidStatus = errors.New(
		"invalid payment status",
	)

	ErrOrderNotPayable = errors.New(
		"order is not payable",
	)
)

type Service struct {
	repo      *Repository
	orderRepo *order.Repository
	providers *ProviderRegistry
}

func NewService(
	repo *Repository,
	orderRepo *order.Repository,
	providers *ProviderRegistry,
) *Service {
	return &Service{
		repo:      repo,
		orderRepo: orderRepo,
		providers: providers,
	}
}

func (s *Service) Initialize(
	ctx context.Context,
	req InitializePaymentRequest,
	userID string,
	email string,
) (*PaymentInitialization, error) {

	orderData, err := s.orderRepo.GetByIDForUser(
		ctx,
		req.OrderID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	if orderData.Status != "pending" {
		return nil, ErrOrderNotPayable
	}

	providerName := strings.ToLower(
		strings.TrimSpace(req.Provider),
	)

	provider, ok := s.providers.Get(providerName)
	if !ok {
		return nil, ErrProviderNotFound
	}

	now := time.Now()

	payment := &Payment{
		ID:        uuid.NewString(),
		OrderID:   orderData.ID,
		UserID:    userID,
		Provider:  providerName,
		Amount:    orderData.Total,
		Currency:  orderData.Currency,
		Status:    PaymentPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	initialization, err := provider.Initialize(
		ctx,
		payment,
		email,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize payment: %w",
			err,
		)
	}

	payment.ProviderRef = initialization.ProviderRef

	if err := s.repo.Create(
		ctx,
		payment,
	); err != nil {
		return nil, fmt.Errorf(
			"create payment: %w",
			err,
		)
	}

	initialization.PaymentID = payment.ID

	return initialization, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (*Payment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Verify(
	ctx context.Context,
	id string,
) (*Payment, error) {

	payment, err := s.repo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return nil, err
	}

	provider, ok := s.providers.Get(
		payment.Provider,
	)
	if !ok {
		return nil, ErrProviderNotFound
	}

	result, err := provider.Verify(
		ctx,
		payment.ProviderRef,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"verify payment: %w",
			err,
		)
	}

	switch result.Status {
	case PaymentPending,
		PaymentSucceeded,
		PaymentFailed,
		PaymentRefunded:

	default:
		return nil, ErrInvalidStatus
	}

	if err := s.repo.UpdateStatus(
		ctx,
		payment.ID,
		result.Status,
	); err != nil {
		return nil, err
	}

	payment.Status = result.Status

	return payment, nil
}
