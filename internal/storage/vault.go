package storage

import (
	"context"
)

type Vault struct {
	ID     uint64 `json:"id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	UserID uint64 `json:"user_id"`
}

//go:generate mockgen -source=vault.go -destination=./mock/vault.go -package=mock
type VaultStorage interface {
	GetVault(ctx context.Context, id uint64) (Vault, error)
	CreateVault(ctx context.Context, v Vault) (Vault, error)
	UpdateVault(ctx context.Context, v Vault) error
	DeleteVault(ctx context.Context, id uint64) error
}
