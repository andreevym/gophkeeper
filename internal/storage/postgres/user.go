package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserStorage handles operations related to user data in a PostgreSQL database.
type UserStorage struct {
	db *sqlx.DB
}

// NewUserStorage creates a new instance of UserStorage.
// It takes a *sqlx.DB instance which is used to interact with the database.
// Returns a pointer to a UserStorage instance.
func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

// GetUser retrieves a user by their ID.
// It takes a context.Context and a user ID (uint64) as parameters.
// Returns a storage.User object and an error if any.
// If the user is not found, it returns ErrUserNotFound.
func (s UserStorage) GetUser(ctx context.Context, id uint64) (storage.User, error) {
	sql := `SELECT login, password FROM users WHERE id = $1`
	u := storage.User{
		ID: id,
	}
	err := s.db.QueryRowContext(ctx, sql, id).Scan(&u.Login, &u.Password)
	if err != nil {
		if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return u, ErrUserNotFound
		}
		return u, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}

	return u, nil
}

// GetUserByLogin retrieves a user by their login.
// It takes a context.Context and a login string as parameters.
// Returns a storage.User object and an error if any.
// If the user is not found, it returns ErrUserNotFound.
func (s UserStorage) GetUserByLogin(ctx context.Context, login string) (storage.User, error) {
	sql := `SELECT id, password FROM users WHERE login = $1`
	u := storage.User{
		Login: login,
	}
	err := s.db.QueryRowContext(ctx, sql, login).Scan(&u.ID, &u.Password)
	if err != nil {
		if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return u, ErrUserNotFound
		}
		return u, fmt.Errorf("failed to get user by login %s: %w", login, err)
	}

	return u, nil
}

// CreateUser inserts a new user into the database.
// It takes a context.Context and a storage.User object as parameters.
// Returns the created storage.User object and an error if any.
func (s UserStorage) CreateUser(ctx context.Context, u storage.User) (storage.User, error) {
	rows, err := s.db.QueryContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", u.Login, u.Password)
	if err != nil {
		return u, fmt.Errorf("failed to create user %s: %w", u.Login, err)
	}

	if rows.Next() {
		err = rows.Scan(&u.ID)
		if err != nil {
			return u, fmt.Errorf("failed to create user %s: %w", u.Login, err)
		}
	}

	return u, nil
}

// UpdateUser updates an existing user in the database.
// It takes a context.Context and a storage.User object as parameters.
// Returns an error if any.
func (s UserStorage) UpdateUser(ctx context.Context, u storage.User) error {
	sql := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := s.db.ExecContext(ctx, sql, u.ID, u.Login, u.Password)
	if err != nil {
		return fmt.Errorf("failed to update user by id %d, login %s: %w", u.ID, u.Login, err)
	}

	return nil
}

// DeleteUser removes a user from the database by their ID.
// It takes a context.Context and a user ID (uint64) as parameters.
// Returns an error if any.
func (s UserStorage) DeleteUser(ctx context.Context, id uint64) error {
	sql := `DELETE FROM users WHERE id = $1`
	_, err := s.db.ExecContext(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete user by id %d: %w", id, err)
	}

	return nil
}
