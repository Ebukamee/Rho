package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateAddress(ctx context.Context, address *Address) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if address.IsDefault {
		_, err = tx.Exec(
			ctx,
			`UPDATE addresses
			 SET is_default = FALSE, updated_at = NOW()
			 WHERE user_id = $1`,
			address.UserID,
		)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO addresses (
			id,
			user_id,
			label,
			first_name,
			last_name,
			company,
			phone,
			line1,
			line2,
			city,
			state,
			postal_code,
			country_code,
			is_default,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14, $15, $16
		)`,
		address.ID,
		address.UserID,
		address.Label,
		address.FirstName,
		address.LastName,
		address.Company,
		address.Phone,
		address.Line1,
		address.Line2,
		address.City,
		address.State,
		address.PostalCode,
		address.CountryCode,
		address.IsDefault,
		address.CreatedAt,
		address.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetAddressByID(
	ctx context.Context,
	userID string,
	addressID string,
) (*Address, error) {
	address := &Address{}

	err := r.db.QueryRow(
		ctx,
		`SELECT
			id,
			user_id,
			label,
			first_name,
			last_name,
			company,
			phone,
			line1,
			line2,
			city,
			state,
			postal_code,
			country_code,
			is_default,
			created_at,
			updated_at
		FROM addresses
		WHERE id = $1 AND user_id = $2`,
		addressID,
		userID,
	).Scan(
		&address.ID,
		&address.UserID,
		&address.Label,
		&address.FirstName,
		&address.LastName,
		&address.Company,
		&address.Phone,
		&address.Line1,
		&address.Line2,
		&address.City,
		&address.State,
		&address.PostalCode,
		&address.CountryCode,
		&address.IsDefault,
		&address.CreatedAt,
		&address.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (r *Repository) ListAddresses(ctx context.Context, userID string) ([]Address, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			id,
			user_id,
			label,
			first_name,
			last_name,
			company,
			phone,
			line1,
			line2,
			city,
			state,
			postal_code,
			country_code,
			is_default,
			created_at,
			updated_at
		FROM addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]Address, 0)

	for rows.Next() {
		var address Address

		err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.Label,
			&address.FirstName,
			&address.LastName,
			&address.Company,
			&address.Phone,
			&address.Line1,
			&address.Line2,
			&address.City,
			&address.State,
			&address.PostalCode,
			&address.CountryCode,
			&address.IsDefault,
			&address.CreatedAt,
			&address.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *Repository) UpdateAddress(ctx context.Context, address *Address) (bool, error) {
	commandTag, err := r.db.Exec(
		ctx,
		`UPDATE addresses
		SET
			label = $1,
			first_name = $2,
			last_name = $3,
			company = $4,
			phone = $5,
			line1 = $6,
			line2 = $7,
			city = $8,
			state = $9,
			postal_code = $10,
			country_code = $11,
			is_default = $12,
			updated_at = $13
		WHERE id = $14 AND user_id = $15`,
		address.Label,
		address.FirstName,
		address.LastName,
		address.Company,
		address.Phone,
		address.Line1,
		address.Line2,
		address.City,
		address.State,
		address.PostalCode,
		address.CountryCode,
		address.IsDefault,
		address.UpdatedAt,
		address.ID,
		address.UserID,
	)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func (r *Repository) SetDefaultAddress(
	ctx context.Context,
	userID string,
	addressID string,
) (bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`UPDATE addresses
		 SET is_default = FALSE, updated_at = NOW()
		 WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return false, err
	}

	commandTag, err := tx.Exec(
		ctx,
		`UPDATE addresses
		 SET is_default = TRUE, updated_at = NOW()
		 WHERE id = $1 AND user_id = $2`,
		addressID,
		userID,
	)
	if err != nil {
		return false, err
	}

	if commandTag.RowsAffected() == 0 {
		return false, pgx.ErrNoRows
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) DeleteAddress(
	ctx context.Context,
	userID string,
	addressID string,
) (bool, error) {
	commandTag, err := r.db.Exec(
		ctx,
		`DELETE FROM addresses
		 WHERE id = $1 AND user_id = $2`,
		addressID,
		userID,
	)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}
