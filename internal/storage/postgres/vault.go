package postgres

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"io"
	"strings"
)

var (
	ErrVaultNotFound = errors.New("vault not found")
)

type VaultStorage struct {
	db   *sqlx.DB
	conn *pgx.Conn
}

func NewVaultStorage(db *sqlx.DB, conn *pgx.Conn) *VaultStorage {
	return &VaultStorage{db: db, conn: conn}
}

func (s VaultStorage) GetVault(ctx context.Context, id uint64) (storage.Vault, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var oid uint32
	var v storage.Vault
	err = s.db.QueryRowContext(ctx, "SELECT key, value, user_id FROM vault WHERE id = $1", id).Scan(&v.Key, &oid, &v.UserID)
	if err != nil {
		if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return v, ErrVaultNotFound
		}
		return v, fmt.Errorf("failed to get vault by id %d: %w", id, err)
	}
	v.ID = id

	lobs := tx.LargeObjects()
	obj, err := lobs.Open(ctx, oid, pgx.LargeObjectModeRead)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to open object: %w", err)
	}

	buffer := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buffer, obj)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to read object: %w", err)
	}
	v.Value = buffer.String()
	return v, nil
}

func (s VaultStorage) CreateVault(ctx context.Context, v storage.Vault) (storage.Vault, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	lobs := tx.LargeObjects()

	oid, err := lobs.Create(ctx, 0)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to create vault object: %w", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO vault (key, value, user_id) VALUES ($1, $2, $3) RETURNING id", v.Key, oid, v.UserID)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to create vault %s: %w", v.Key, err)
	}

	// Open the new Object for writing.
	obj, err := lobs.Open(ctx, oid, pgx.LargeObjectModeWrite)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to open vault %s: %w", v.Key, err)
	}

	// Copy the file stream to the Large Object stream
	_, err = io.Copy(obj, bytes.NewReader([]byte(v.Value)))
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to copy vault %s: %w", v.Key, err)
	}

	err = tx.Commit(ctx)
	return storage.Vault{}, err
}

func (s VaultStorage) UpdateVault(ctx context.Context, v storage.Vault) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	lobs := tx.LargeObjects()

	oid, err := lobs.Create(ctx, 0)
	if err != nil {
		return fmt.Errorf("failed to create vault object: %w", err)
	}

	_, err = s.db.ExecContext(ctx, "UPDATE vault SET key = $2, value = $3 WHERE id = $1", v.ID, v.Key, oid)
	if err != nil {
		return fmt.Errorf("failed to update vault by id %d, key %s: %w", v.ID, v.Key, err)
	}

	// Open the new Object for writing.
	obj, err := lobs.Open(ctx, oid, pgx.LargeObjectModeWrite)
	if err != nil {
		return fmt.Errorf("failed to open vault %s: %w", v.Key, err)
	}

	// Copy the file stream to the Large Object stream
	_, err = io.Copy(obj, bytes.NewReader([]byte(v.Value)))
	if err != nil {
		return fmt.Errorf("failed to copy vault %s: %w", v.Key, err)
	}

	err = tx.Commit(ctx)
	return nil
}

func (s VaultStorage) DeleteVault(ctx context.Context, id uint64) error {
	sql := `DELETE FROM vault WHERE id = $1`
	_, err := s.db.ExecContext(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete vault by id %d: %w", id, err)
	}

	return nil
}
