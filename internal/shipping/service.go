package shipping

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req CreateShipmentRequest,
	userID string,
) (*Shipment, error) {
	now := time.Now()

	shipment := &Shipment{
		ID:             uuid.NewString(),
		OrderID:        strings.TrimSpace(req.OrderID),
		UserID:         userID,
		Carrier:        strings.TrimSpace(req.Carrier),
		Service:        strings.TrimSpace(req.Service),
		TrackingNumber: strings.TrimSpace(req.TrackingNumber),
		Cost:           req.Cost,
		Currency:       strings.ToUpper(strings.TrimSpace(req.Currency)),
		Status:         StatusPending,
		EstimatedDays:  req.EstimatedDays,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.Create(ctx, shipment); err != nil {
		return nil, fmt.Errorf(
			"create shipment: %w",
			err,
		)
	}

	return shipment, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (*Shipment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByOrderID(
	ctx context.Context,
	orderID string,
) (*Shipment, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	req UpdateShipmentRequest,
) (*Shipment, error) {
	shipment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Carrier != "" {
		shipment.Carrier = strings.TrimSpace(req.Carrier)
	}

	if req.Service != "" {
		shipment.Service = strings.TrimSpace(req.Service)
	}

	if req.TrackingNumber != "" {
		shipment.TrackingNumber =
			strings.TrimSpace(req.TrackingNumber)
	}

	if req.Status != "" {
		shipment.Status = req.Status
	}

	if req.EstimatedDays != nil {
		shipment.EstimatedDays = *req.EstimatedDays
	}

	shipment.UpdatedAt = time.Now()

	if err := s.repo.Update(
		ctx,
		shipment,
	); err != nil {
		return nil, fmt.Errorf(
			"update shipment: %w",
			err,
		)
	}

	return shipment, nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf(
			"delete shipment: %w",
			err,
		)
	}

	return nil
}