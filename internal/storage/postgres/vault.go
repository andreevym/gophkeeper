package postgres

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

var (
	ErrVaultNotFound = errors.New("vault not found")
)

// VaultStorage handles operations related to vault data in a PostgreSQL database.
type VaultStorage struct {
	db   *sqlx.DB
	conn *pgx.Conn
}

// NewVaultStorage creates a new instance of VaultStorage.
// It takes a *sqlx.DB and a *pgx.Conn instance which are used to interact with the database.
// Returns a pointer to a VaultStorage instance.
func NewVaultStorage(db *sqlx.DB, conn *pgx.Conn) *VaultStorage {
	return &VaultStorage{db: db, conn: conn}
}

// GetVault retrieves a vault by its ID.
// It takes a context.Context and a vault ID (uint64) as parameters.
// Returns a storage.Vault object and an error if any.
// If the vault is not found, it returns ErrVaultNotFound.
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

// CreateVault inserts a new vault into the database.
// It takes a context.Context and a storage.Vault object as parameters.
// Returns the created storage.Vault object and an error if any.
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

	var vaultID uint64

	err = tx.QueryRow(ctx, "INSERT INTO vault (key, value, user_id) VALUES ($1, $2, $3) RETURNING id", v.Key, oid, v.UserID).Scan(&vaultID)
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
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return storage.Vault{
		ID:     vaultID,
		Key:    v.Key,
		Value:  v.Value,
		UserID: v.UserID,
	}, nil
}

// UpdateVault updates an existing vault in the database.
// It takes a context.Context and a storage.Vault object as parameters.
// Returns an error if any.
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

	_, err = tx.Exec(ctx, "UPDATE vault SET key = $2, value = $3 WHERE id = $1", v.ID, v.Key, oid)
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
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteVault removes a vault from the database by its ID.
// It takes a context.Context and a vault ID (uint64) as parameters.
// Returns an error if any.
func (s VaultStorage) DeleteVault(ctx context.Context, id uint64) error {
	sql := `DELETE FROM vault WHERE id = $1`
	_, err := s.db.ExecContext(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete vault by id %d: %w", id, err)
	}

	return nil
}
