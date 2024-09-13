package handlers

import "net/http"

func (h *ServiceHandlers) GetPingHandler(writer http.ResponseWriter, request *http.Request) {
	if h.dbClient == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx := request.Context()
	err := h.dbClient.PingContext(ctx)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
