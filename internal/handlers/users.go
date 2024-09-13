package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/andreevym/gophkeeper/internal/storage"
	storage2 "github.com/andreevym/gophkeeper/internal/storage/postgres"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
}

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *ServiceHandlers) PostSignUp(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Logger().Warn("failed to read all bytes", zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	signUpRequest := SignUpRequest{}
	err = json.Unmarshal(bytes, &signUpRequest)
	if err != nil {
		logger.Logger().Warn("failed to unmarshal post signup request", zap.String("signUpRequest", string(bytes)), zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	_, err = h.userStorage.GetUserByLogin(ctx, signUpRequest.Login)
	if err != nil && !errors.Is(storage2.ErrUserNotFound, err) {
		logger.Logger().Warn("failed to get user by login", zap.String("login", signUpRequest.Login), zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), bcrypt.MaxCost)
	if err != nil {
		logger.Logger().Warn("failed to generate hash from password", zap.String("login", signUpRequest.Login), zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	user := storage.User{
		Login:    signUpRequest.Login,
		Password: string(hashedPassword),
	}
	_, err = h.userStorage.CreateUser(ctx, user)
	if err != nil {
		logger.Logger().Warn("failed to create user", zap.String("login", signUpRequest.Login), zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}
}

func (h *ServiceHandlers) PostSignIn(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Logger().Warn("failed to read all bytes", zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	signInRequest := SignInRequest{}
	err = json.Unmarshal(bytes, &signInRequest)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal post sign in request %s: %w", signInRequest.Login, err)
		logger.Logger().Warn("failed to unmarshal post sign in request", zap.String("signInRequest", string(bytes)), zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	user, err := h.userStorage.GetUserByLogin(ctx, signInRequest.Login)
	if err != nil {
		err = fmt.Errorf("failed to get user by login %s: %w", signInRequest.Login, err)
		logger.Logger().Warn("failed to get user by login", zap.String("login", signInRequest.Login), zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	if !h.hasher.Match(user.Password, signInRequest.Password) {
		msg := fmt.Sprintf("failed to match password %s", signInRequest.Login)
		logger.Logger().Warn(msg, zap.String("login", signInRequest.Login))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, msg)
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	authToken, err := h.authProvider.GenerateToken(user.ID)
	if err != nil {
		err = fmt.Errorf("failed to generate token %s: %w", signInRequest.Login, err)
		logger.Logger().Warn("failed to generate token", zap.String("login", signInRequest.Login), zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	writer.WriteHeader(http.StatusOK)
}
