package postgres_test

import (
	"context"
	"log"
	"testing"

	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/andreevym/gophkeeper/internal/storage/postgres"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	db := postgres.NewDB()
	defer db.TeardownDB()
	err := db.SetupDB("../../../migrations")
	if err != nil {
		log.Fatalf("Could not setup postgres container: %v", err)
	}

	ctx := context.Background()
	err = db.DB.PingContext(ctx)
	require.NoError(t, err)

	userStorage := postgres.NewUserStorage(db.DB)

	user1 := storage.User{
		Login:    "k1",
		Password: "v1",
	}

	_, err = userStorage.CreateUser(ctx, user1)
	require.NoError(t, err)

	user2 := storage.User{
		Login:    "k2",
		Password: "v2",
	}

	_, err = userStorage.CreateUser(ctx, user2)
	require.NoError(t, err)

	founduser1, err := userStorage.GetUser(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, founduser1)
	require.Equal(t, user1.Login, founduser1.Login)
	require.Equal(t, user1.Password, founduser1.Password)

	updateduser1 := storage.User{
		ID:       1,
		Login:    "k1",
		Password: "k3",
	}

	err = userStorage.UpdateUser(ctx, updateduser1)
	require.NoError(t, err)

	afterUpdateuser1, err := userStorage.GetUser(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, afterUpdateuser1)
	require.Equal(t, updateduser1.Login, afterUpdateuser1.Login)
	require.Equal(t, updateduser1.Password, afterUpdateuser1.Password)

	err = userStorage.DeleteUser(ctx, 1)
	require.NoError(t, err)

	_, err = userStorage.GetUser(ctx, 1)
	require.EqualError(t, err, postgres.ErrUserNotFound.Error())
}
