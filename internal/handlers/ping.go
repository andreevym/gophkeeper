package handlers

import "net/http"

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
