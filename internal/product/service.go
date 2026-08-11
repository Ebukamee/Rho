package product

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rho-commerce/rho/pkg/pagination"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*Product, error) {
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))

	product := &Product{
		ID:          uuid.NewString(),
		Name:        strings.TrimSpace(req.Name),
		Slug:        strings.TrimSpace(req.Slug),
		Description: req.Description,
		SKU:         strings.TrimSpace(req.SKU),
		Price:       req.Price,
		Currency:    currency,
		ImageURL:    req.ImageURL,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.Active != nil {
		product.Active = *req.Active
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return product, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, page, limit int, activeOnly bool) (*ProductListResponse, error) {
	products, total, err := s.repo.List(ctx, page, limit, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	return &ProductListResponse{
		Products:   products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: pagination.TotalPages(total, limit),
	}, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateProductRequest) (*Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = strings.TrimSpace(req.Name)
	}

	if req.Slug != "" {
		product.Slug = strings.TrimSpace(req.Slug)
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.SKU != "" {
		product.SKU = strings.TrimSpace(req.SKU)
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.Currency != "" {
		product.Currency = strings.ToUpper(strings.TrimSpace(req.Currency))
	}

	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}

	if req.Active != nil {
		product.Active = *req.Active
	}

	product.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return product, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	return nil
}
