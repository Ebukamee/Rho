package inventory

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rho-commerce/rho/internal/database"
)

var (
	ErrNotFound = errors.New("inventory not found")
	ErrConflict = errors.New("inventory already exists")
)

type Repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *Repository {
	return &Repository{db: db}
}

func NewRepositoryWithTx(tx pgx.Tx) *Repository {
	return &Repository{db: tx}
}

func (r *Repository) Create(
	ctx context.Context,
	inventory *Inventory,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO inventory (
			id,
			product_id,
			quantity,
			reserved,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6)
	`,
		inventory.ID,
		inventory.ProductID,
		inventory.Quantity,
		inventory.Reserved,
		inventory.CreatedAt,
		inventory.UpdatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}

	return err
}

func (r *Repository) GetByProductID(
	ctx context.Context,
	productID string,
) (*Inventory, error) {
	var inventory Inventory

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			product_id,
			quantity,
			reserved,
			created_at,
			updated_at
		FROM inventory
		WHERE product_id = $1
	`, productID).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Reserved,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*Inventory, error) {
	var inventory Inventory

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			product_id,
			quantity,
			reserved,
			created_at,
			updated_at
		FROM inventory
		WHERE id = $1
	`, id).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Reserved,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *Repository) Update(
	ctx context.Context,
	inventory *Inventory,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE inventory
		SET
			quantity = $1,
			reserved = $2,
			updated_at = $3
		WHERE id = $4
	`,
		inventory.Quantity,
		inventory.Reserved,
		inventory.UpdatedAt,
		inventory.ID,
	)

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
		DELETE FROM inventory
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

func (r *Repository) Adjust(
	ctx context.Context,
	productID string,
	quantity int,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE inventory
		SET
			quantity = quantity + $1,
			updated_at = NOW()
		WHERE product_id = $2
		AND quantity + $1 >= reserved
	`,
		quantity,
		productID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Reserve(
	ctx context.Context,
	productID string,
	quantity int,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE inventory
		SET
			reserved = reserved + $1,
			updated_at = NOW()
		WHERE product_id = $2
		AND quantity - reserved >= $1
	`,
		quantity,
		productID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Release(
	ctx context.Context,
	productID string,
	quantity int,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE inventory
		SET
			reserved = reserved - $1,
			updated_at = NOW()
		WHERE product_id = $2
		AND reserved >= $1
	`,
		quantity,
		productID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
