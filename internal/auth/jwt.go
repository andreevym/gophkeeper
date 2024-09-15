package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/andreevym/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

func GenKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Ошибка при генерации ключа:", err)
		return nil, err
	}

	return privateKey, nil
}

func GenPrivateKeyMust() *ecdsa.PrivateKey {
	key, err := GenKey()
	if err != nil {
		panic(err)
	}
	return key
}

func ReadJwtSecretKey(jwtSecretKey string) (*ecdsa.PrivateKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(jwtSecretKey)
	if err != nil {
		logger.Logger().Error("failed to decode 'jwtSecretKey'", zap.Error(err))
		return nil, fmt.Errorf("failed to decode 'jwtSecretKey': %w", err)
	}
	privateKey, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		logger.Logger().Error("failed to parse ec private key", zap.Error(err))
		return nil, fmt.Errorf("failed to parse ec private key: %w", err)
	}
	return privateKey, nil
}

func MakeJwtSecretKey() (*ecdsa.PrivateKey, string, error) {
	jwtSecretKey := GenPrivateKeyMust()
	privateKey, err := x509.MarshalECPrivateKey(jwtSecretKey)
	if err != nil {
		logger.Logger().Error("failed to marshal private key", zap.Error(err))
		return nil, "", fmt.Errorf("failed to marshal private key: %w", err)
	}
	return jwtSecretKey, base64.StdEncoding.EncodeToString(privateKey), nil
}
