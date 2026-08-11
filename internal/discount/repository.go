package discount

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("discount not found")
	ErrConflict = errors.New("discount code already exists")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	discount *Discount,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO discounts (
			id,
			code,
			type,
			value,
			minimum_order,
			usage_limit,
			usage_count,
			starts_at,
			expires_at,
			active,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
		)
	`,
		discount.ID,
		discount.Code,
		discount.Type,
		discount.Value,
		discount.MinimumOrder,
		discount.UsageLimit,
		discount.UsageCount,
		discount.StartsAt,
		discount.ExpiresAt,
		discount.Active,
		discount.CreatedAt,
		discount.UpdatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}

	return err
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*Discount, error) {
	var d Discount

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			code,
			type,
			value,
			minimum_order,
			usage_limit,
			usage_count,
			starts_at,
			expires_at,
			active,
			created_at,
			updated_at
		FROM discounts
		WHERE id = $1
	`, id).Scan(
		&d.ID,
		&d.Code,
		&d.Type,
		&d.Value,
		&d.MinimumOrder,
		&d.UsageLimit,
		&d.UsageCount,
		&d.StartsAt,
		&d.ExpiresAt,
		&d.Active,
		&d.CreatedAt,
		&d.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *Repository) GetByCode(
	ctx context.Context,
	code string,
) (*Discount, error) {
	var d Discount

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			code,
			type,
			value,
			minimum_order,
			usage_limit,
			usage_count,
			starts_at,
			expires_at,
			active,
			created_at,
			updated_at
		FROM discounts
		WHERE code = $1
	`, code).Scan(
		&d.ID,
		&d.Code,
		&d.Type,
		&d.Value,
		&d.MinimumOrder,
		&d.UsageLimit,
		&d.UsageCount,
		&d.StartsAt,
		&d.ExpiresAt,
		&d.Active,
		&d.CreatedAt,
		&d.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *Repository) Update(
	ctx context.Context,
	d *Discount,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE discounts
		SET
			code = $1,
			type = $2,
			value = $3,
			minimum_order = $4,
			usage_limit = $5,
			starts_at = $6,
			expires_at = $7,
			active = $8,
			updated_at = $9
		WHERE id = $10
	`,
		d.Code,
		d.Type,
		d.Value,
		d.MinimumOrder,
		d.UsageLimit,
		d.StartsAt,
		d.ExpiresAt,
		d.Active,
		d.UpdatedAt,
		d.ID,
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

func (r *Repository) Delete(
	ctx context.Context,
	id string,
) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM discounts
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) IncrementUsage(
	ctx context.Context,
	id string,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE discounts
		SET usage_count = usage_count + 1,
		    updated_at = NOW()
		WHERE id = $1
		  AND (
		      usage_limit IS NULL
		      OR usage_count < usage_limit
		  )
	`, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("discount usage limit reached")
	}

	return nil
}
