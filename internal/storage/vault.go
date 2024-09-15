package storage

import (
	"context"
)

// Vault represents a vault entity with an ID, key, value, and associated user ID.
type Vault struct {
	ID     uint64 `json:"id"`      // Unique identifier for the vault.
	Key    string `json:"key"`     // The key or name of the vault.
	Value  []byte `json:"value"`   // The value or content stored in the vault.
	UserID uint64 `json:"user_id"` // The ID of the user who owns the vault.
}

// VaultStorage defines the interface for operations on vault entities in the storage system.
type VaultStorage interface {
	// GetVault retrieves a vault by its unique ID.
	// Takes a context.Context and the vault's ID (uint64) as parameters.
	// Returns the Vault and an error if any.
	GetVault(ctx context.Context, id uint64) (Vault, error)

	// CreateVault inserts a new vault into the storage system.
	// Takes a context.Context and a Vault object as parameters.
	// Returns the created Vault and an error if any.
	CreateVault(ctx context.Context, v Vault) (Vault, error)

	// UpdateVault updates an existing vault's information.
	// Takes a context.Context and a Vault object with updated information as parameters.
	// Returns an error if any.
	UpdateVault(ctx context.Context, v Vault) error

	// DeleteVault removes a vault from the storage system by its ID.
	// Takes a context.Context and the vault's ID (uint64) as parameters.
	// Returns an error if any.
	DeleteVault(ctx context.Context, id uint64) error
}
