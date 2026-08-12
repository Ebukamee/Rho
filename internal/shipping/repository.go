package shipping

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("shipment not found")
	ErrConflict = errors.New("shipment already exists")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	shipment *Shipment,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO shipments (
			id,
			order_id,
			user_id,
			carrier,
			service,
			tracking_number,
			cost,
			currency,
			status,
			estimated_days,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,
			$7,$8,$9,$10,$11,$12
		)
	`,
		shipment.ID,
		shipment.OrderID,
		shipment.UserID,
		shipment.Carrier,
		shipment.Service,
		shipment.TrackingNumber,
		shipment.Cost,
		shipment.Currency,
		shipment.Status,
		shipment.EstimatedDays,
		shipment.CreatedAt,
		shipment.UpdatedAt,
	)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) &&
		pgErr.Code == "23505" {
		return ErrConflict
	}

	return err
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*Shipment, error) {
	var shipment Shipment

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			order_id,
			user_id,
			carrier,
			service,
			tracking_number,
			cost,
			currency,
			status,
			estimated_days,
			created_at,
			updated_at
		FROM shipments
		WHERE id = $1
	`, id).Scan(
		&shipment.ID,
		&shipment.OrderID,
		&shipment.UserID,
		&shipment.Carrier,
		&shipment.Service,
		&shipment.TrackingNumber,
		&shipment.Cost,
		&shipment.Currency,
		&shipment.Status,
		&shipment.EstimatedDays,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *Repository) GetByOrderID(
	ctx context.Context,
	orderID string,
) (*Shipment, error) {
	var shipment Shipment

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			order_id,
			user_id,
			carrier,
			service,
			tracking_number,
			cost,
			currency,
			status,
			estimated_days,
			created_at,
			updated_at
		FROM shipments
		WHERE order_id = $1
	`, orderID).Scan(
		&shipment.ID,
		&shipment.OrderID,
		&shipment.UserID,
		&shipment.Carrier,
		&shipment.Service,
		&shipment.TrackingNumber,
		&shipment.Cost,
		&shipment.Currency,
		&shipment.Status,
		&shipment.EstimatedDays,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *Repository) Update(
	ctx context.Context,
	shipment *Shipment,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE shipments
		SET
			carrier = $1,
			service = $2,
			tracking_number = $3,
			status = $4,
			estimated_days = $5,
			updated_at = $6
		WHERE id = $7
	`,
		shipment.Carrier,
		shipment.Service,
		shipment.TrackingNumber,
		shipment.Status,
		shipment.EstimatedDays,
		shipment.UpdatedAt,
		shipment.ID,
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
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM shipments WHERE id = $1`,
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