package postgres_test

import (
	"context"
	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/andreevym/gophkeeper/internal/storage/postgres"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestVaultRepository(t *testing.T) {
	db := postgres.NewDB()
	defer db.TeardownDB()
	err := db.SetupDB("../../migrations")
	if err != nil {
		log.Fatalf("Could not setup postgres container: %v", err)
	}

	ctx := context.Background()
	err = db.DB.PingContext(ctx)
	require.NoError(t, err)

	userStorage := postgres.NewUserStorage(db.DB)
	u, err := userStorage.CreateUser(ctx, storage.User{
		Login:    "test",
		Password: "test",
	})
	require.NoError(t, err)

	vaultStorage := postgres.NewVaultStorage(db.DB, db.Conn)

	vault1 := storage.Vault{
		Key:    "k1",
		Value:  "v1",
		UserID: u.ID,
	}

	_, err = vaultStorage.CreateVault(ctx, vault1)
	require.NoError(t, err)

	vault2 := storage.Vault{
		Key:    "k1",
		Value:  "v1",
		UserID: u.ID,
	}

	_, err = vaultStorage.CreateVault(ctx, vault2)
	require.NoError(t, err)

	foundVault1, err := vaultStorage.GetVault(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, foundVault1)
	require.Equal(t, vault1.Key, foundVault1.Key)
	require.Equal(t, vault1.Value, foundVault1.Value)
	require.Equal(t, vault1.UserID, foundVault1.UserID)

	updatedVault1 := storage.Vault{
		ID:     1,
		Key:    "k1",
		Value:  "v2",
		UserID: u.ID,
	}

	err = vaultStorage.UpdateVault(ctx, updatedVault1)
	require.NoError(t, err)

	afterUpdateVault1, err := vaultStorage.GetVault(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, afterUpdateVault1)
	require.Equal(t, updatedVault1.Key, afterUpdateVault1.Key)
	require.Equal(t, updatedVault1.Value, afterUpdateVault1.Value)
	require.Equal(t, updatedVault1.UserID, afterUpdateVault1.UserID)

	err = vaultStorage.DeleteVault(ctx, 1)
	require.NoError(t, err)

	_, err = vaultStorage.GetVault(ctx, 1)
	require.EqualError(t, err, postgres.ErrVaultNotFound.Error())
}
