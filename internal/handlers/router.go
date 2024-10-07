package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Hasher defines methods for hashing and checking passwords.
type Hasher interface {
	Hash(password string) (string, error)       // Hash generates a hash of the password.
	Match(hashedPassword, password string) bool // Match checks if the provided password matches the hashed password.
}

// UserSessionExtractor defines methods for handling user sessions and JWT tokens.
type UserSessionExtractor interface {
	GenerateToken(userID uint64) (string, error)                  // GenerateToken generates a JWT token for the given user ID.
	GetUserFromSession(ctx context.Context) (storage.User, error) // GetUserFromSession retrieves the user from the session context.
}

// Constants for various URI paths used in the application.
const (
	RootURI       = "/"                // RootURI is the root endpoint.
	AuthSignInURI = "/api/auth/signin" // AuthSignInURI is the endpoint for user sign-in.
	AuthSignUpURI = "/api/auth/signup" // AuthSignUpURI is the endpoint for user sign-up.
	VaultURI      = "/api/vault"       // VaultURI is the endpoint for vault operations.
	PingURI       = "/api/ping"        // PingURI is the endpoint for health checks.
	FileUploadURI = "/api/upload"      // FileUploadURI is the endpoint for upload file.
)

// ServiceHandlers manages HTTP request handlers for the service.
type ServiceHandlers struct {
	dbClient     DBClient             // DBClient interface for interacting with the database.
	authProvider UserSessionExtractor // UserSessionExtractor for handling user sessions and tokens.
	userStorage  storage.UserStorage  // UserStorage for user-related operations.
	vaultStorage storage.VaultStorage // VaultStorage for vault-related operations.
	hashService  Hasher               // Hasher for password hashing and verification.
}

// DBClient defines methods for database operations.
type DBClient interface {
	PingContext(ctx context.Context) error // PingContext checks the database connection.
}

// NewServiceHandlers creates a new instance of ServiceHandlers with the given dependencies.
// It initializes the ServiceHandlers with the provided DBClient, UserSessionExtractor,
// VaultStorage, UserStorage, and Hasher.
//
// Parameters:
//   - dbClient (DBClient): Interface for database operations.
//   - authProvider (UserSessionExtractor): Interface for user sessions and JWT tokens.
//   - vaultStorage (storage.VaultStorage): Interface for vault operations.
//   - userStorage (storage.UserStorage): Interface for user operations.
//   - hashService (Hasher): Interface for password hashing.
//
// Returns:
//   - *ServiceHandlers: A new instance of ServiceHandlers with the provided dependencies.
func NewServiceHandlers(
	dbClient DBClient,
	authProvider UserSessionExtractor,
	vaultStorage storage.VaultStorage,
	userStorage storage.UserStorage,
	hashService Hasher,
) *ServiceHandlers {
	return &ServiceHandlers{
		dbClient:     dbClient,
		authProvider: authProvider,
		vaultStorage: vaultStorage,
		userStorage:  userStorage,
		hashService:  hashService,
	}
}

// NewRouter creates a new HTTP router with the specified handlers and middleware.
//
// Parameters:
//   - s (*ServiceHandlers): The service handlers containing the business logic.
//   - m (...func(http.Handler) http.Handler): Optional middleware functions to apply to the router.
//
// Returns:
//   - *chi.Mux: A new router instance with the specified routes and middleware.
func NewRouter(s *ServiceHandlers, m ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(m...)

	// Route handlers
	r.Post(AuthSignInURI, s.PostSignIn)
	r.Post(AuthSignUpURI, s.PostSignUp)

	r.Post(VaultURI, s.PostVault)
	r.Get(VaultURI+"/{vaultID}", s.GetVault)

	r.Get(PingURI, s.GetPingHandler)

	r.Post(FileUploadURI, s.FileUploadHandler)
	r.Post(FileUploadURI+"/{vaultID}", s.FileUploadHandler)

	r.Get(RootURI, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
	})

	return r
}
