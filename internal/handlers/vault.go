package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// VaultRequest represents the payload for vault operations.
type VaultRequest struct {
	ID    string `json:"id"`    // The ID of the vault entry. Empty for new entries.
	Key   string `json:"key"`   // The key for the vault entry.
	Value string `json:"value"` // The value for the vault entry.
}

// PostVault handles both the creation and update of vault entries.
//
// If the ID is empty, a new vault entry is created with the provided Key and Value.
// If the ID is provided, the vault entry is updated if it exists and belongs to the current user.
//
// The handler responds with:
//   - HTTP 400 Bad Request if there are errors in request processing or validation.
//   - HTTP 201 Created if a new vault entry is created successfully.
//   - HTTP 200 OK if an existing vault entry is updated successfully.
func (h *ServiceHandlers) PostVault(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.authProvider.GetUserFromSession(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to validate user session: %v", err), http.StatusBadRequest)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vaultRequest := VaultRequest{}
	err = json.Unmarshal(bytes, &vaultRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		bytes, err := json.Marshal(vault)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(bytes)
		if err != nil {
			logger.Logger().Warn("failed to write response", zap.Error(err))
		}
		return
	}

	id, err := strconv.ParseUint(vaultRequest.ID, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse param vaultID: %v", err), http.StatusBadRequest)
		return
	}

	v, err := h.vaultStorage.GetVault(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get vault: %v", err), http.StatusBadRequest)
		return
	}

	if v.UserID != user.ID {
		http.Error(w, "access denied", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err = json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		logger.Logger().Warn("failed to write response", zap.Error(err))
	}
	w.WriteHeader(http.StatusOK)
}

// GetVault handles the retrieval of a specific vault entry.
//
// The vaultID is extracted from the request URL path. If the vault entry exists and belongs to the current user,
// it is returned in the response.
//
// The handler responds with:
//   - HTTP 400 Bad Request if there are errors in request processing or validation.
//   - HTTP 200 OK with the vault entry details if successfully retrieved.
func (h *ServiceHandlers) GetVault(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	user, err := h.authProvider.GetUserFromSession(ctx)
	if err != nil {
		http.Error(writer, fmt.Sprintf("failed to validate user session: %v", err), http.StatusBadRequest)
		return
	}

	vaultID := chi.URLParam(request, "vaultID")
	if vaultID == "" {
		http.Error(writer, "vaultID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(vaultID, 10, 64)
	if err != nil {
		http.Error(writer, fmt.Sprintf("failed to parse param vaultID: %v", err), http.StatusBadRequest)
		return
	}

	v, err := h.vaultStorage.GetVault(ctx, id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("failed to get vault: %v", err), http.StatusBadRequest)
		return
	}

	if v.UserID != user.ID {
		http.Error(writer, "access denied", http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
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
