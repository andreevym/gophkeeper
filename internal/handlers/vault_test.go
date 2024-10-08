package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreevym/gophkeeper/internal/auth"
	"github.com/andreevym/gophkeeper/internal/handlers"
	"github.com/andreevym/gophkeeper/internal/middleware"
	"github.com/andreevym/gophkeeper/internal/pwd"
	"github.com/andreevym/gophkeeper/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultHandler(t *testing.T) {
	t.Parallel()
	db := postgres.NewDB()
	defer db.TeardownDB()

	ctx := context.Background()
	err := db.SetupDB(ctx, "../../migrations")
	if err != nil {
		log.Fatalf("Could not setup postgres container: %v", err)
	}

	_, jwtSecretKey, err := auth.MakeJwtSecretKey()
	require.NoError(t, err)

	vaultStorage := postgres.NewVaultStorage(db.DB, db.Conn)
	userStorage := postgres.NewUserStorage(db.DB)

	jwtPrivateKey, err := auth.ReadJwtSecretKey(jwtSecretKey)
	require.NoError(t, err)

	authProvider := auth.NewAuthProvider(userStorage, jwtPrivateKey)

	authMiddleware := auth.NewAuthMiddleware(authProvider, jwtSecretKey, handlers.AuthSignInURI, handlers.AuthSignUpURI)

	hashService := pwd.NewHashService()
	serviceHandlers := handlers.NewServiceHandlers(db.DB, authProvider, vaultStorage, userStorage, hashService)

	router := handlers.NewRouter(
		serviceHandlers,
		authMiddleware.WithAuthentication,
		middleware.WithRequestLoggerMiddleware,
	)
	ts := httptest.NewServer(router)
	defer ts.Close()
	signUpRequest := handlers.SignUpRequest{
		Login:    "test",
		Password: "test",
	}
	reqBody, err := json.Marshal(signUpRequest)
	require.NoError(t, err)
	header := http.Header{}
	statusCode, _, got := testRequest(t, ts, http.MethodPost, handlers.AuthSignUpURI, bytes.NewBuffer(reqBody), header)
	require.Equal(t, http.StatusCreated, statusCode, "failed to sign up user", string(reqBody), got)
	assert.Equal(t, "{\"id\":1,\"login\":\"test\"}", got)

	signInRequest := handlers.SignInRequest{
		Login:    signUpRequest.Login,
		Password: signUpRequest.Password,
	}
	reqBody, err = json.Marshal(signInRequest)
	require.NoError(t, err)
	statusCode, header, got = testRequest(t, ts, http.MethodPost, handlers.AuthSignInURI, bytes.NewBuffer(reqBody), header)
	require.Equal(t, http.StatusOK, statusCode, "failed to sign in user", string(reqBody), got)
	assert.Empty(t, got)

	vaultRequest := handlers.VaultRequest{
		Key:   "key",
		Value: "val",
	}
	reqBody, err = json.Marshal(vaultRequest)
	require.NoError(t, err)
	statusCode, _, got = testRequest(t, ts, http.MethodPost, handlers.VaultURI, bytes.NewBuffer(reqBody), header)
	require.Equal(t, http.StatusCreated, statusCode, "failed to make vault request", string(reqBody), got)
	assert.Contains(t, got, fmt.Sprintf("{\"id\":1,\"key\":\"key\",\"value\":\"val\",\"user_id\":"))

	require.NotEmpty(t, got)

	vaultResponse := handlers.VaultResponse{}
	err = json.Unmarshal([]byte(got), &vaultResponse)
	require.NoError(t, err)

	statusCode, _, got = testRequest(t, ts, http.MethodGet, fmt.Sprintf("%s/%d", handlers.VaultURI, vaultResponse.ID), bytes.NewBuffer(reqBody), header)
	require.Equal(t, http.StatusOK, statusCode)
	assert.Contains(t, got, "{\"id\":1,\"key\":\"key\",\"value\":\"val\",\"user_id\":")
}
