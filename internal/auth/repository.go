package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, role, provider, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.Role, user.Provider, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password, first_name, last_name, role, provider, created_at, updated_at
		FROM users WHERE email = $1`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.Provider, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, email, password, first_name, last_name, role, provider, created_at, updated_at
		FROM users WHERE id = $1`

	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.Provider, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) ListUsers(ctx context.Context, page, limit int) ([]User, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, email, password, first_name, last_name, role, provider, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName,
			&u.Role, &u.Provider, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}
func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	query := `
		UPDATE users 
		SET email = $1, password = $2, first_name = $3, last_name = $4, 
		    role = $5, provider = $6, updated_at = $7
		WHERE id = $8`

	_, err := r.db.Exec(ctx, query,
		user.Email, user.Password, user.FirstName, user.LastName,
		user.Role, user.Provider, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *Repository) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
