package auth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"strconv"
)

type Provider struct {
	userStorage   storage.UserStorage
	jwtPrivateKey *ecdsa.PrivateKey
}

func NewAuthProvider(userStorage storage.UserStorage, jwtPrivateKey *ecdsa.PrivateKey) *Provider {
	return &Provider{
		userStorage:   userStorage,
		jwtPrivateKey: jwtPrivateKey,
	}
}

type ContextKey int

var ErrAuthUnauthorized = errors.New("unauthorized")

const (
	UserIDContextKey ContextKey = iota
)

func (p Provider) CreateSession(ctx context.Context, userID uint64) (context.Context, error) {
	_, err := p.userStorage.GetUser(ctx, userID)
	if err != nil {
		return ctx, fmt.Errorf("get user: %w", err)
	}

	ctxWithValue := context.WithValue(ctx, UserIDContextKey, userID)
	return ctxWithValue, nil
}

func (p Provider) GetUserFromSession(ctx context.Context) (storage.User, error) {
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
