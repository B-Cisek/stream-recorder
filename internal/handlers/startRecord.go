package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/B-Cisek/stream-recorder/internal/services"
)

type StartRecordHandler struct {
	service services.RecorderService
}

func NewStartRecordHandler(service services.RecorderService) *StartRecordHandler {
	return &StartRecordHandler{
		service: service,
	}
}

type StartRecordRequest struct {
	Channel  string `json:"channel"`
	Platform string `json:"platform"`
}

type StartRecordResponse struct {
	Message string `json:"message"`
}

func (h *StartRecordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.StartRecording(req.Channel); err != nil {
		http.Error(w, "Failed to start recording", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := StartRecordResponse{
		Message: "recording started for channel: " + req.Channel,
	}

	json.NewEncoder(w).Encode(response)
}
