package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/andreevym/gophkeeper/internal/storage"
	storage2 "github.com/andreevym/gophkeeper/internal/storage/postgres"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// SignUpRequest represents the payload for user sign-up requests.
type SignUpRequest struct {
	Login    string `json:"login"`    // The login username for the new user.
	Password string `json:"password"` // The password for the new user.
}

// SignUpResponse represents the response payload for user sign-up requests.
type SignUpResponse struct {
	ID    uint64 `json:"id"`
	Login string `json:"login"`
}

// SignInRequest represents the payload for user sign-in requests.
type SignInRequest struct {
	Login    string `json:"login"`    // The login username for the user.
	Password string `json:"password"` // The password for the user.
}

// PostSignUp handles user sign-up requests.
//
// It reads the user credentials from the request body, validates them,
// hashes the password, and stores the new user in the database.
//
// The handler responds with:
//   - HTTP 400 Bad Request if there is an error in the request or processing.
//   - HTTP 201 Created on successful user creation.
func (h *ServiceHandlers) PostSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Logger().Warn("failed to read all bytes", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signUpRequest := SignUpRequest{}
	err = json.Unmarshal(bytes, &signUpRequest)
	if err != nil {
		logger.Logger().Warn("failed to unmarshal post signup request", zap.String("signUpRequest", string(bytes)), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if signUpRequest.Login == "" || len(signUpRequest.Login) > 50 {
		logger.Logger().Warn("login is empty or too long more than 50 characters", zap.Int("LoginLen", len(signUpRequest.Login)))
		http.Error(w, fmt.Sprintf("login is empty or too long more than 50 characters but actual len is %d", len(signUpRequest.Login)), http.StatusBadRequest)
		return
	}

	if signUpRequest.Password == "" || len(signUpRequest.Password) > 50 {
		logger.Logger().Warn("password is empty or too long more than 50 characters", zap.Int("PasswordLen", len(signUpRequest.Password)))
		http.Error(w, fmt.Sprintf("password is empty or too long more than 50 characters but actual len is %d", len(signUpRequest.Password)), http.StatusBadRequest)
		return
	}

	_, err = h.userStorage.GetUserByLogin(ctx, signUpRequest.Login)
	if err != nil && !errors.Is(storage2.ErrUserNotFound, err) {
		logger.Logger().Warn("failed to get user by login", zap.String("login", signUpRequest.Login), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err == nil {
		logger.Logger().Info("user already exists", zap.String("login", signUpRequest.Login))
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), bcrypt.MinCost)
	if err != nil {
		logger.Logger().Warn("failed to generate hash from password", zap.String("login", signUpRequest.Login), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := storage.User{
		Login:    signUpRequest.Login,
		Password: string(hashedPassword),
	}
	createdUser, err := h.userStorage.CreateUser(ctx, user)
	if err != nil {
		logger.Logger().Warn("failed to create user", zap.String("login", signUpRequest.Login), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := SignUpResponse{
		ID:    createdUser.ID,
		Login: createdUser.Login,
	}

	bytes, err = json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(bytes)
	if err != nil {
		logger.Logger().Warn("failed to write response", zap.Error(err))
	}
}

// PostSignIn handles user sign-in requests.
//
// It reads the user credentials from the request body, validates the user,
// and generates a JWT token if the credentials are correct.
//
// The handler responds with:
//   - HTTP 400 Bad Request if there is an error in the request or credentials are invalid.
//   - HTTP 200 OK with the JWT token in the Authorization header on successful sign-in.
func (h *ServiceHandlers) PostSignIn(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Logger().Warn("failed to read all bytes", zap.Error(err))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	signInRequest := SignInRequest{}
	err = json.Unmarshal(bytes, &signInRequest)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal post sign in request %s: %w", signInRequest.Login, err)
		logger.Logger().Warn("failed to unmarshal post sign in request", zap.String("signInRequest", string(bytes)), zap.Error(err))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userStorage.GetUserByLogin(ctx, signInRequest.Login)
	if err != nil {
		err = fmt.Errorf("failed to get user by login %s: %w", signInRequest.Login, err)
		logger.Logger().Warn("failed to get user by login", zap.String("login", signInRequest.Login), zap.Error(err))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.hashService.Match(user.Password, signInRequest.Password) {
		msg := fmt.Sprintf("failed to match password %s", signInRequest.Login)
		logger.Logger().Warn(msg, zap.String("login", signInRequest.Login))
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	authToken, err := h.authProvider.GenerateToken(user.ID)
	if err != nil {
		err = fmt.Errorf("failed to generate token %s: %w", signInRequest.Login, err)
		logger.Logger().Warn("failed to generate token", zap.String("login", signInRequest.Login), zap.Error(err))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	writer.WriteHeader(http.StatusOK)
}
