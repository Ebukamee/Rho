package category

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

func (s *Service) Create(
	ctx context.Context,
	req CreateCategoryRequest,
) (*Category, error) {
	category := &Category{
		ID:          uuid.NewString(),
		Name:        strings.TrimSpace(req.Name),
		Slug:        strings.TrimSpace(req.Slug),
		Description: strings.TrimSpace(req.Description),
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.Active != nil {
		category.Active = *req.Active
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	return category, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (*Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(
	ctx context.Context,
	page int,
	limit int,
	activeOnly bool,
) (*CategoryListResponse, error) {
	categories, total, err := s.repo.List(
		ctx,
		page,
		limit,
		activeOnly,
	)

	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return &CategoryListResponse{
		Categories: categories,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: pagination.TotalPages(total, limit),
	}, nil
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	req UpdateCategoryRequest,
) (*Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = strings.TrimSpace(req.Name)
	}

	if req.Slug != "" {
		category.Slug = strings.TrimSpace(req.Slug)
	}

	if req.Description != "" {
		category.Description = strings.TrimSpace(req.Description)
	}

	if req.Active != nil {
		category.Active = *req.Active
	}

	category.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	return category, nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}
