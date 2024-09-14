package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/andreevym/gophkeeper/internal/auth"
	"github.com/andreevym/gophkeeper/internal/handlers"
	"github.com/andreevym/gophkeeper/internal/middleware"
	"github.com/andreevym/gophkeeper/internal/pwd"
	"github.com/andreevym/gophkeeper/internal/storage/postgres"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler(t *testing.T) {
	db := postgres.NewDB()
	defer db.TeardownDB()
	err := db.SetupDB("../../migrations")
	if err != nil {
		log.Fatalf("Could not setup postgres container: %v", err)
	}

	type want struct {
		resp       string
		statusCode int
	}
	tests := []struct {
		name          string
		want          want
		request       string
		httpMethod    string
		signUpRequest *handlers.SignUpRequest
	}{
		{
			name: "success register new user",
			want: want{
				statusCode: http.StatusOK,
				resp:       "",
			},
			request: handlers.AuthSignUpURI,
			signUpRequest: &handlers.SignUpRequest{
				Login:    "a",
				Password: "b",
			},
			httpMethod: http.MethodPost,
		},
		{
			name: "success register new user (50 char password and login)",
			want: want{
				statusCode: http.StatusOK,
				resp:       "",
			},
			request: handlers.AuthSignUpURI,
			signUpRequest: &handlers.SignUpRequest{
				Login:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Password: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			},
			httpMethod: http.MethodPost,
		},
		{
			name: "success register new user (empty password and login)",
			want: want{
				statusCode: http.StatusBadRequest,
				resp:       "login is empty or too long more than 50 characters but actual len is 0",
			},
			request: handlers.AuthSignUpURI,
			signUpRequest: &handlers.SignUpRequest{
				Login:    "",
				Password: "",
			},
			httpMethod: http.MethodPost,
		},
		{
			name: "success register new user (more 50 char password and login)",
			want: want{
				statusCode: http.StatusBadRequest,
				resp:       "login is empty or too long more than 50 characters but actual len is 51",
			},
			request: handlers.AuthSignUpURI,
			signUpRequest: &handlers.SignUpRequest{
				Login:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Password: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			},
			httpMethod: http.MethodPost,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
			var reqBody []byte
			if test.signUpRequest != nil {
				reqBody, err = json.Marshal(test.signUpRequest)
				require.NoError(t, err)
			} else {
				reqBody = []byte{}
			}
			header := http.Header{}
			statusCode, _, got := testRequest(t, ts, test.httpMethod, test.request, bytes.NewBuffer(reqBody), header)
			assert.Equal(t, test.want.statusCode, statusCode)
			assert.Equal(t, test.want.resp, got)
		})
	}
}
