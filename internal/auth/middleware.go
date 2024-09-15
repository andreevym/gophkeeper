package auth

import (
	"fmt"
	"net/http"

	"github.com/andreevym/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

// Middleware provides HTTP middleware for authentication using JWT tokens.
type Middleware struct {
	authProvider         *Provider           // Provider for authentication and JWT handling
	jwtSecretKey         string              // JWT secret key used for token validation
	allowUnauthorizedURI map[string]struct{} // URIs that can be accessed without authentication
}

// JwtService defines an interface for JWT token validation.
type JwtService interface {
	ValidateToken(tokenString string, jwtSecretKey string) (uint64, error)
}

// NewAuthMiddleware creates a new instance of Middleware with the given AuthProvider and JWT secret key.
//
// Parameters:
//   - authProvider (*Provider): The provider used for authentication and JWT handling.
//   - jwtSecretKey (string): The JWT secret key used for validating tokens.
//   - allowUnauthorizedURI (...string): List of URIs that can be accessed without authentication.
//
// Returns:
//   - *Middleware: A new instance of Middleware configured with the provided parameters.
func NewAuthMiddleware(
	authProvider *Provider,
	jwtSecretKey string,
	allowUnauthorizedURI ...string,
) *Middleware {
	allowUnauthorizedURIMap := make(map[string]struct{})
	for _, uri := range allowUnauthorizedURI {
		allowUnauthorizedURIMap[uri] = struct{}{}
	}

	return &Middleware{
		authProvider:         authProvider,
		jwtSecretKey:         jwtSecretKey,
		allowUnauthorizedURI: allowUnauthorizedURIMap,
	}
}

// WithAuthentication returns an HTTP handler that performs authentication based on JWT tokens.
// Requests to URIs in `allowUnauthorizedURI` are allowed without authentication.
//
// Parameters:
//   - next (http.Handler): The next handler to call if authentication is successful.
//
// Returns:
//   - http.Handler: An HTTP handler with authentication middleware applied.
func (m *Middleware) WithAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request URI is allowed without authentication
		if _, ok := m.allowUnauthorizedURI[r.RequestURI]; ok {
			next.ServeHTTP(w, r)
			return
		}

		// Extract the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix to get the token string
		tokenString := authHeader[len("Bearer "):]

		// Validate the token and extract user ID
		userID, err := m.authProvider.ValidateToken(tokenString)
		if err != nil {
			logger.Logger().Warn("jwtService.ValidateToken", zap.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set the user ID from the token in the request context
		ctx, err := m.authProvider.CreateSession(r.Context(), userID)
		if err != nil {
			logger.Logger().Warn("create session", zap.Error(err))
			http.Error(w, fmt.Sprintf("failed to create session: %v", err), http.StatusUnauthorized)
			return
		}

		// Pass the context with the user ID to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
