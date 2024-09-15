package handlers

import (
	"context"
	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type Hasher interface {
	Hash(password string) (string, error)
	Match(hashedPassword, password string) bool
}

type UserSessionExtractor interface {
	GenerateToken(userID uint64) (string, error)
	GetUserFromSession(ctx context.Context) (storage.User, error)
}

const (
	RootURI       = "/"
	AuthSignInURI = "/api/auth/signin"
	AuthSignUpURI = "/api/auth/signup"
	VaultURI      = "/api/vault"
	PingURI       = "/api/ping"
)

type ServiceHandlers struct {
	dbClient     *sqlx.DB
	authProvider UserSessionExtractor
	userStorage  storage.UserStorage
	vaultStorage storage.VaultStorage
	hasher       Hasher
}

func NewServiceHandlers(
	dbClient *sqlx.DB,
	authProvider UserSessionExtractor,
	vaultStorage storage.VaultStorage,
	userStorage storage.UserStorage,
	hasher Hasher,
) *ServiceHandlers {
	return &ServiceHandlers{
		dbClient:     dbClient,
		authProvider: authProvider,
		vaultStorage: vaultStorage,
		userStorage:  userStorage,
		hasher:       hasher,
	}
}

// NewRouter creates a new HTTP router with the specified handlers and tracer.
func NewRouter(s *ServiceHandlers, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewares...)

	r.Post(AuthSignInURI, s.PostSignIn)
	r.Post(AuthSignUpURI, s.PostSignUp)

	r.Post(VaultURI, s.PostVault)
	r.Get(VaultURI+"/{vaultID}", s.GetVault)

	r.Get(PingURI, s.GetPingHandler)

	r.Get(RootURI, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
	})
	return r
}
