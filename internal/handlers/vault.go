package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type VaultRequest struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *ServiceHandlers) PostVault(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	user, err := h.authProvider.GetUserFromSession(ctx)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, fmt.Sprintf("failed to validate user session: %v", err))
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	vaultRequest := VaultRequest{}
	err = json.Unmarshal(bytes, &vaultRequest)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	if vaultRequest.ID == "" {
		vault := storage.Vault{
			Key:    vaultRequest.Key,
			Value:  vaultRequest.Value,
			UserID: user.ID,
		}
		vault, err = h.vaultStorage.CreateVault(ctx, vault)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, err := io.WriteString(writer, err.Error())
			if err != nil {
				logger.Logger().Warn("failed to write response", zap.Error(err))
			}
			return
		}

		bytes, err := json.Marshal(vault)
		if err != nil {
			_, err := io.WriteString(writer, err.Error())
			if err != nil {
				logger.Logger().Warn("failed to write response", zap.Error(err))
			}
			return
		}
		_, err = writer.Write(bytes)
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		writer.WriteHeader(http.StatusCreated)
		return
	}

	id, err := strconv.ParseUint(vaultRequest.ID, 10, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, fmt.Sprintf("failed to parse param vaultID: %v", err))
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	v, err := h.vaultStorage.GetVault(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, fmt.Sprintf("failed to get vault: %v", err))
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	if v.UserID != user.ID {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, "access denied")
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	if vaultRequest.Key != "" {
		v.Key = vaultRequest.Key
	}

	if vaultRequest.Value != "" {
		v.Value = vaultRequest.Value
	}

	err = h.vaultStorage.UpdateVault(ctx, v)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	bytes, err = json.Marshal(v)
	if err != nil {
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}
	_, err = writer.Write(bytes)
	if err != nil {
		logger.Logger().Warn("failed to write response", zap.Error(err))
	}
	writer.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) GetVault(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	user, err := h.authProvider.GetUserFromSession(ctx)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, fmt.Sprintf("failed to validate user session: %v", err))
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	vaultID := chi.URLParam(request, "vaultID")
	if vaultID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, "vaultID is required")
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}
	id, err := strconv.ParseUint(vaultID, 10, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, fmt.Sprintf("failed to parse param vaultID: %v", err))
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	v, err := h.vaultStorage.GetVault(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, fmt.Sprintf("failed to get vault: %v", err))
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	if v.UserID != user.ID {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, "access denied")
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		_, err := io.WriteString(writer, err.Error())
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}
	_, err = writer.Write(bytes)
	if err != nil {
		logger.Logger().Warn("failed to write response", zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
