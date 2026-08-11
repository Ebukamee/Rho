package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("order not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	order *Order,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (
			id,
			user_id,
			status,
			subtotal,
			discount,
			total,
			currency,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9
		)
	`,
		order.ID,
		order.UserID,
		order.Status,
		order.Subtotal,
		order.Discount,
		order.Total,
		order.Currency,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (
				id,
				order_id,
				product_id,
				name,
				sku,
				quantity,
				unit_price,
				total_price
			)
			VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8
			)
		`,
			item.ID,
			order.ID,
			item.ProductID,
			item.Name,
			item.SKU,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*Order, error) {

	var order Order

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			status,
			subtotal,
			discount,
			total,
			currency,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
	`, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Subtotal,
		&order.Discount,
		&order.Total,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			order_id,
			product_id,
			name,
			sku,
			quantity,
			unit_price,
			total_price
		FROM order_items
		WHERE order_id = $1
	`, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order.Items = []OrderItem{}

	for rows.Next() {
		var item OrderItem

		rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Name,
			&item.SKU,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
		)

		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func (r *Repository) UpdateStatus(
	ctx context.Context,
	id string,
	status OrderStatus,
) error {

	result, err := r.db.Exec(ctx, `
		UPDATE orders
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
	`, status, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
