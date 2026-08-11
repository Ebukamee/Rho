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
			id, category_id, name, slug, description, sku, price,
			currency, image_url, active, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`

	_, err := r.db.Exec(ctx, query,
		product.ID,
		product.CategoryID,
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
		SELECT
			id, category_id, name, slug, description, sku, price,
			currency, image_url, active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var p Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.CategoryID,
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

func (r *Repository) List(
	ctx context.Context,
	page, limit int,
	activeOnly bool,
	categoryID *string,
) ([]Product, int, error) {
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM products WHERE 1=1`
	countArgs := []any{}

	if activeOnly {
		countQuery += ` AND active = TRUE`
	}

	if categoryID != nil {
		countQuery += ` AND category_id = $1`
		countArgs = append(countArgs, *categoryID)
	}

	var total int

	if err := r.db.QueryRow(
		ctx,
		countQuery,
		countArgs...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id, category_id, name, slug, description, sku, price,
			currency, image_url, active, created_at, updated_at
		FROM products
		WHERE 1=1
	`

	args := []any{}
	argPosition := 1

	if activeOnly {
		query += ` AND active = TRUE`
	}

	if categoryID != nil {
		query += ` AND category_id = $` + string(rune('0'+argPosition))
		args = append(args, *categoryID)
		argPosition++
	}

	query += ` ORDER BY created_at DESC`
	query += ` LIMIT $` + string(rune('0'+argPosition))
	args = append(args, limit)
	argPosition++

	query += ` OFFSET $` + string(rune('0'+argPosition))
	args = append(args, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var p Product

		if err := rows.Scan(
			&p.ID,
			&p.CategoryID,
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
		SET category_id = $1,
		    name = $2,
		    slug = $3,
		    description = $4,
		    sku = $5,
		    price = $6,
		    currency = $7,
		    image_url = $8,
		    active = $9,
		    updated_at = $10
		WHERE id = $11
	`

	result, err := r.db.Exec(ctx, query,
		product.CategoryID,
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
	result, err := r.db.Exec(
		ctx,
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
