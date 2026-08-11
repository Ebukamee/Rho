package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("product not found")
	ErrConflict = errors.New("product slug or SKU already exists")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, product *Product) error {
	query := `
		INSERT INTO products (
			id, name, slug, description, sku, price,
			currency, image_url, active, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err := r.db.Exec(ctx, query,
		product.ID,
		product.Name,
		product.Slug,
		product.Description,
		product.SKU,
		product.Price,
		product.Currency,
		product.ImageURL,
		product.Active,
		product.CreatedAt,
		product.UpdatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}

	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Product, error) {
	query := `
		SELECT id, name, slug, description, sku, price,
		       currency, image_url, active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var p Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Slug,
		&p.Description,
		&p.SKU,
		&p.Price,
		&p.Currency,
		&p.ImageURL,
		&p.Active,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *Repository) List(ctx context.Context, page, limit int, activeOnly bool) ([]Product, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM products`

	if activeOnly {
		countQuery += ` WHERE active = TRUE`
	}

	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, slug, description, sku, price,
		       currency, image_url, active, created_at, updated_at
		FROM products
	`

	if activeOnly {
		query += ` WHERE active = TRUE`
	}

	query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var p Product

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.SKU,
			&p.Price,
			&p.Currency,
			&p.ImageURL,
			&p.Active,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *Repository) Update(ctx context.Context, product *Product) error {
	query := `
		UPDATE products
		SET name = $1,
		    slug = $2,
		    description = $3,
		    sku = $4,
		    price = $5,
		    currency = $6,
		    image_url = $7,
		    active = $8,
		    updated_at = $9
		WHERE id = $10
	`

	result, err := r.db.Exec(ctx, query,
		product.Name,
		product.Slug,
		product.Description,
		product.SKU,
		product.Price,
		product.Currency,
		product.ImageURL,
		product.Active,
		product.UpdatedAt,
		product.ID,
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
	result, err := r.db.Exec(ctx,
		`DELETE FROM products WHERE id = $1`,
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
