package handlers

import (
	"encoding/json"
	"net/http"
)

type PingHandler struct {
}

type PingResponse struct {
	Message string `json:"message"`
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(PingResponse{
		Message: "ok",
	})
}
