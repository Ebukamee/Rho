package payment

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("payment not found")
	ErrConflict = errors.New("payment already exists")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	payment *Payment,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO payments (
			id,
			order_id,
			user_id,
			provider,
			provider_ref,
			amount,
			currency,
			status,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10
		)
	`,
		payment.ID,
		payment.OrderID,
		payment.UserID,
		payment.Provider,
		payment.ProviderRef,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.CreatedAt,
		payment.UpdatedAt,
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
) (*Payment, error) {
	var p Payment

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			order_id,
			user_id,
			provider,
			provider_ref,
			amount,
			currency,
			status,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
	`, id).Scan(
		&p.ID,
		&p.OrderID,
		&p.UserID,
		&p.Provider,
		&p.ProviderRef,
		&p.Amount,
		&p.Currency,
		&p.Status,
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

func (r *Repository) GetByProviderRef(
	ctx context.Context,
	providerRef string,
) (*Payment, error) {
	var p Payment

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			order_id,
			user_id,
			provider,
			provider_ref,
			amount,
			currency,
			status,
			created_at,
			updated_at
		FROM payments
		WHERE provider_ref = $1
	`, providerRef).Scan(
		&p.ID,
		&p.OrderID,
		&p.UserID,
		&p.Provider,
		&p.ProviderRef,
		&p.Amount,
		&p.Currency,
		&p.Status,
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

func (r *Repository) UpdateStatus(
	ctx context.Context,
	id string,
	status PaymentStatus,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE payments
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
	`,
		status,
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
