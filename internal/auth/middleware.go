package auth

import (
	"fmt"
	"net/http"

	"github.com/andreevym/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

// Middleware is a middleware for authentication using JWT tokens.
type Middleware struct {
	authProvider         *Provider
	jwtSecretKey         string
	allowUnauthorizedURI map[string]struct{}
}

type JwtService interface {
	ValidateToken(tokenString string, jwtSecretKey string) (uint64, error)
}

// NewAuthMiddleware creates a new instance of Middleware with the given AuthService.
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

// WithAuthentication implements the http.HandlerFunc interface for the Middleware.
func (m *Middleware) WithAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
