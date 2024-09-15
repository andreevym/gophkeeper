package handlers

import (
	"net/http"
)

// GetPingHandler handles the HTTP GET request for the /ping endpoint.
// It checks the database connection by calling PingContext on the dbClient.
// If the database client is not initialized or if PingContext returns an error,
// it responds with an HTTP 500 Internal Server Error status.
// Otherwise, it responds with an HTTP 200 OK status.
//
// Parameters:
//   - w (http.ResponseWriter): The response writer to send the response.
//   - r (http.Request): The HTTP request containing the context.
//
// Responses:
//   - HTTP 200 OK: The database connection is healthy.
//   - HTTP 500 Internal Server Error: The database connection is not healthy or dbClient is nil.
func (h *ServiceHandlers) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	if h.dbClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	err := h.dbClient.PingContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
