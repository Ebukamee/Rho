package payment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProviderNotFound = errors.New("payment provider not found")
	ErrInvalidStatus    = errors.New("invalid payment status")
)

type Service struct {
	repo      *Repository
	providers *ProviderRegistry
}

func NewService(
	repo *Repository,
	providers *ProviderRegistry,
) *Service {
	return &Service{
		repo:      repo,
		providers: providers,
	}
}

func (s *Service) Initialize(
	ctx context.Context,
	req InitializePaymentRequest,
	orderID string,
	userID string,
	amount int64,
	currency string,
) (*PaymentInitialization, error) {

	providerName := strings.ToLower(strings.TrimSpace(req.Provider))

	provider, ok := s.providers.Get(providerName)
	if !ok {
		return nil, ErrProviderNotFound
	}

	now := time.Now()

	payment := &Payment{
		ID:        uuid.NewString(),
		OrderID:   orderID,
		UserID:    userID,
		Provider:  providerName,
		Amount:    amount,
		Currency:  currency,
		Status:    PaymentPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	initialization, err := provider.Initialize(
		ctx,
		payment,
	)
	if err != nil {
		return nil, fmt.Errorf("initialize payment: %w", err)
	}

	payment.ProviderRef = initialization.ProviderRef

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
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

	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	provider, ok := s.providers.Get(payment.Provider)
	if !ok {
		return nil, ErrProviderNotFound
	}

	result, err := provider.Verify(
		ctx,
		payment.ProviderRef,
	)
	if err != nil {
		return nil, fmt.Errorf("verify payment: %w", err)
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
