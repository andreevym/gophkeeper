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

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

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

func (s UserStorage) UpdateUser(ctx context.Context, u storage.User) error {
	sql := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := s.db.ExecContext(ctx, sql, u.ID, u.Login, u.Password)
	if err != nil {
		return fmt.Errorf("failed to update user by id %d, login %s: %w", u.ID, u.Login, err)
	}

	return nil
}

func (s UserStorage) DeleteUser(ctx context.Context, id uint64) error {
	sql := `DELETE FROM users WHERE id = $1`
	_, err := s.db.ExecContext(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete user by id %d: %w", id, err)
	}

	return nil
}
