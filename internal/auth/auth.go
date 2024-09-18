package auth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strconv"

	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Provider is a structure that handles authentication and authorization.
// It uses a storage system for user data and an ECDSA private key for JWT signing and verification.
type Provider struct {
	userStorage   storage.UserStorage // Storage interface for user data operations
	jwtPrivateKey *ecdsa.PrivateKey   // ECDSA private key for signing JWTs
}

// NewAuthProvider creates a new instance of Provider with the given user storage and JWT private key.
//
// Parameters:
//   - userStorage (storage.UserStorage): The storage interface for user data.
//   - jwtPrivateKey (*ecdsa.PrivateKey): The private key used for signing JWTs.
//
// Returns:
//   - *Provider: A new Provider instance.
func NewAuthProvider(userStorage storage.UserStorage, jwtPrivateKey *ecdsa.PrivateKey) *Provider {
	return &Provider{
		userStorage:   userStorage,
		jwtPrivateKey: jwtPrivateKey,
	}
}

// ContextKey is a custom type for context keys used in the authentication process.
type ContextKey int

const (
	// UserIDContextKey is the context key used to store and retrieve user IDs in the context.
	UserIDContextKey ContextKey = iota
)

// ErrAuthUnauthorized is an error returned when a user is unauthorized.
var ErrAuthUnauthorized = errors.New("unauthorized")

// CreateSession creates a new session for the user with the specified userID and stores it in the context.
//
// Parameters:
//   - ctx (context.Context): The context in which the session will be created.
//   - userID (uint64): The ID of the user for whom the session is being created.
//
// Returns:
//   - context.Context: The context with the user ID added.
//   - error: An error if the user cannot be retrieved or if there is an issue creating the session.
func (p *Provider) CreateSession(ctx context.Context, userID uint64) (context.Context, error) {
	_, err := p.userStorage.GetUser(ctx, userID)
	if err != nil {
		return ctx, fmt.Errorf("get user: %w", err)
	}

	ctxWithValue := context.WithValue(ctx, UserIDContextKey, userID)
	return ctxWithValue, nil
}

// GetUserFromSession retrieves the user associated with the session from the context.
//
// Parameters:
//   - ctx (context.Context): The context containing the session information.
//
// Returns:
//   - storage.User: The user associated with the session.
//   - error: An error if the user cannot be retrieved or if there is an issue accessing the session.
func (p *Provider) GetUserFromSession(ctx context.Context) (storage.User, error) {
	userID := ctx.Value(UserIDContextKey)
	if userID == nil {
		return storage.User{}, ErrAuthUnauthorized
	}

	user, err := p.userStorage.GetUser(ctx, userID.(uint64))
	if err != nil {
		return storage.User{}, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

// ValidateToken validates a JWT token and extracts the user ID from it.
//
// Parameters:
//   - tokenString (string): The JWT token to be validated.
//
// Returns:
//   - uint64: The user ID extracted from the token, or 0 if the token is invalid.
//   - error: An error if the token is invalid or if there is an issue parsing the token.
func (p *Provider) ValidateToken(tokenString string) (uint64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return &p.jwtPrivateKey.PublicKey, nil
	})
	notFoundID := uint64(0)
	if err != nil {
		return notFoundID, fmt.Errorf("jwt parse: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return notFoundID, errors.New("token is not valid")
	}

	// Extract the user ID from the token claims
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return notFoundID, errors.New("invalid token claims")
	}
	id := mapClaims["userID"].(string)
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return notFoundID, fmt.Errorf("strconv.ParseInt, %s: invalid user ID in token: %w", id, err)
	}

	return userID, nil
}

// GenerateToken generates a JWT token for a given user ID.
//
// Parameters:
//   - userID (uint64): The ID of the user for whom the token is being generated.
//
// Returns:
//   - string: The generated JWT token.
//   - error: An error if token generation fails.
func (p *Provider) GenerateToken(userID uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, &jwt.MapClaims{
		"userID": strconv.FormatUint(userID, 10),
	})

	t, err := token.SignedString(p.jwtPrivateKey)
	if err != nil {
		logger.Logger().Error("failed to sign token", zap.Error(err))
		return "", fmt.Errorf("sign the token: %w", err)
	}
	return t, nil
}
