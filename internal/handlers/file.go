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

// FileUploadHandler handles the uploading of binary files and manages vault entries.
//
// Depending on whether a vault ID is provided, this handler either creates a new vault entry or
// updates an existing one. The uploaded file will be saved to the server's storage system. The handler
// processes the file from the request's multipart form data and performs the appropriate action based on
// the presence of a vault ID.
func (h *ServiceHandlers) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Retrieve and validate the user session.
	user, err := h.authProvider.GetUserFromSession(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to validate user session: %v", err), http.StatusUnauthorized)
		return
	}

	// Limit the size of the uploaded file to 10MB.
	const maxUploadSize = 10 << 20 // 10MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse the form to retrieve file information.
	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form.
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get file from form: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content.
	bytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read file from form: %v", err), http.StatusInternalServerError)
		return
	}

	vaultID := chi.URLParam(r, "vaultID")

	var v storage.Vault
	if vaultID == "" {
		// Create a new vault entry if no ID is provided.
		v, err = h.vaultStorage.CreateVault(ctx, storage.Vault{
			Key:    header.Filename,
			Value:  bytes,
			UserID: user.ID,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create vault: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		id, err := strconv.ParseUint(vaultID, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse param vaultID '%s': %v", vaultID, err), http.StatusBadRequest)
			return
		}
		// Update an existing vault entry if an ID is provided.
		v = storage.Vault{
			ID:     id,
			Key:    header.Filename,
			Value:  bytes,
			UserID: user.ID,
		}
		err = h.vaultStorage.UpdateVault(ctx, v)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update vault: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	// Marshal the response to JSON.
	resp, err := json.Marshal(VaultResponse{
		ID:     v.ID,
		Key:    v.Key,
		UserID: v.UserID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		logger.Logger().Warn("failed to write response", zap.Error(err))
	}
}
