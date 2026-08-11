package category

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("category not found")
	ErrConflict = errors.New("category slug already exists")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, category *Category) error {
	query := `
		INSERT INTO categories (
			id,
			name,
			slug,
			description,
			active,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		category.ID,
		category.Name,
		category.Slug,
		category.Description,
		category.Active,
		category.CreatedAt,
		category.UpdatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}

	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Category, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			active,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1
	`

	var category Category

	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.Active,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *Repository) List(
	ctx context.Context,
	page int,
	limit int,
	activeOnly bool,
) ([]Category, int, error) {
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM categories`

	if activeOnly {
		countQuery += ` WHERE active = TRUE`
	}

	var total int

	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id,
			name,
			slug,
			description,
			active,
			created_at,
			updated_at
		FROM categories
	`

	if activeOnly {
		query += ` WHERE active = TRUE`
	}

	query += `
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	categories := make([]Category, 0)

	for rows.Next() {
		var category Category

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Active,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *Repository) Update(ctx context.Context, category *Category) error {
	query := `
		UPDATE categories
		SET
			name = $1,
			slug = $2,
			description = $3,
			active = $4,
			updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.Active,
		category.UpdatedAt,
		category.ID,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM categories WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
