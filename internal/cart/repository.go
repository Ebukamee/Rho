package cart

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rho-commerce/rho/internal/database"
)

var (
	ErrCartNotFound     = errors.New("cart not found")
	ErrCartItemNotFound = errors.New("cart item not found")
	ErrDuplicateItem    = errors.New("product already in cart")
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

func (r *Repository) GetByUserID(ctx context.Context, userID string) (*Cart, error) {
	cart := &Cart{}

	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, created_at, updated_at
		FROM carts
		WHERE user_id = $1
	`, userID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCartNotFound
	}

	if err != nil {
		return nil, err
	}

	cart.Items, err = r.getItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *Repository) Create(ctx context.Context, cart *Cart) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO carts (
			id, user_id, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4)
	`,
		cart.ID,
		cart.UserID,
		cart.CreatedAt,
		cart.UpdatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrDuplicateItem
	}

	return err
}

func (r *Repository) GetOrCreate(ctx context.Context, cart *Cart) (*Cart, error) {
	existing, err := r.GetByUserID(ctx, cart.UserID)

	if err == nil {
		return existing, nil
	}

	if !errors.Is(err, ErrCartNotFound) {
		return nil, err
	}

	if err := r.Create(ctx, cart); err != nil {
		if errors.Is(err, ErrDuplicateItem) {
			return r.GetByUserID(ctx, cart.UserID)
		}

		return nil, err
	}

	cart.Items = []CartItem{}

	return cart, nil
}

func (r *Repository) AddItem(ctx context.Context, item *CartItem) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO cart_items (
			id, cart_id, product_id, quantity, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		item.ID,
		item.CartID,
		item.ProductID,
		item.Quantity,
		item.CreatedAt,
		item.UpdatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrDuplicateItem
	}

	return err
}

func (r *Repository) GetItem(
	ctx context.Context,
	cartID string,
	itemID string,
) (*CartItem, error) {
	item := &CartItem{}

	err := r.db.QueryRow(ctx, `
		SELECT id, cart_id, product_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE id = $1 AND cart_id = $2
	`, itemID, cartID).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCartItemNotFound
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *Repository) UpdateItem(
	ctx context.Context,
	cartID string,
	itemID string,
	quantity int,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE cart_items
		SET quantity = $1,
		    updated_at = NOW()
		WHERE id = $2 AND cart_id = $3
	`,
		quantity,
		itemID,
		cartID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCartItemNotFound
	}

	return nil
}

func (r *Repository) RemoveItem(
	ctx context.Context,
	cartID string,
	itemID string,
) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM cart_items
		WHERE id = $1 AND cart_id = $2
	`,
		itemID,
		cartID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCartItemNotFound
	}

	return nil
}

func (r *Repository) Clear(ctx context.Context, cartID string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM cart_items
		WHERE cart_id = $1
	`, cartID)

	return err
}

func (r *Repository) getItems(
	ctx context.Context,
	cartID string,
) ([]CartItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, cart_id, product_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE cart_id = $1
		ORDER BY created_at ASC
	`, cartID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]CartItem, 0)

	for rows.Next() {
		var item CartItem

		if err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
