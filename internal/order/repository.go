package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rho-commerce/rho/internal/database"
)

var ErrNotFound = errors.New("order not found")

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
	order *Order,
) error {
	_, err := r.db.Exec(ctx, `
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
		_, err := r.db.Exec(ctx, `
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

	return nil
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

	if err != nil {
		return nil, err
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

	order.Items = make([]OrderItem, 0)

	for rows.Next() {
		var item OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Name,
			&item.SKU,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
		); err != nil {
			return nil, err
		}

		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

func (r *Repository) GetByIDForUser(
	ctx context.Context,
	orderID string,
	userID string,
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
		WHERE id = $1 AND user_id = $2
	`, orderID, userID).Scan(
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

	if err != nil {
		return nil, err
	}

	return &order, nil
}
