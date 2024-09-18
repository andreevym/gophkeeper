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

// GenKey generates a new ECDSA private key using the P-256 curve.
// It uses a cryptographically secure random number generator to create the key.
//
// Returns:
//   - *ecdsa.PrivateKey: The generated ECDSA private key.
//   - error: An error if key generation fails.
func GenKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Ошибка при генерации ключа:", err)
		return nil, err
	}

	return privateKey, nil
}

// GenPrivateKeyMust generates a new ECDSA private key and panics if an error occurs.
// This function is typically used in scenarios where the application should not continue
// if key generation fails.
//
// Returns:
//   - *ecdsa.PrivateKey: The generated ECDSA private key.
func GenPrivateKeyMust() *ecdsa.PrivateKey {
	key, err := GenKey()
	if err != nil {
		panic(err)
	}
	return key
}

// ReadJwtSecretKey decodes a Base64-encoded JWT secret key and parses it into an ECDSA private key.
// The key is expected to be in PEM format.
//
// Parameters:
//   - jwtSecretKey (string): The Base64-encoded JWT secret key.
//
// Returns:
//   - *ecdsa.PrivateKey: The parsed ECDSA private key.
//   - error: An error if decoding or parsing fails.
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

// MakeJwtSecretKey generates a new ECDSA private key and returns both the key and its
// Base64-encoded representation. This can be used to create a new JWT secret key.
//
// Returns:
//   - *ecdsa.PrivateKey: The generated ECDSA private key.
//   - string: The Base64-encoded representation of the private key.
//   - error: An error if key generation or encoding fails.
func MakeJwtSecretKey() (*ecdsa.PrivateKey, string, error) {
	jwtSecretKey := GenPrivateKeyMust()
	privateKey, err := x509.MarshalECPrivateKey(jwtSecretKey)
	if err != nil {
		logger.Logger().Error("failed to marshal private key", zap.Error(err))
		return nil, "", fmt.Errorf("failed to marshal private key: %w", err)
	}
	return jwtSecretKey, base64.StdEncoding.EncodeToString(privateKey), nil
}
