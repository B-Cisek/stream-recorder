package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/B-Cisek/stream-recorder/internal/services"
)

type StopRecordHandler struct {
	service services.RecorderService
}

func NewStopRecordHandler(service services.RecorderService) *StopRecordHandler {
	return &StopRecordHandler{
		service: service,
	}
}

type StopRecordRequest struct {
	Channel  string `json:"channel"`
	Platform string `json:"platform"`
}

type StopRecordResponse struct {
	Message string `json:"message"`
}

func (h *StopRecordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StopRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.StopRecording(req.Platform, req.Channel); err != nil {
		http.Error(w, "Failed to stop recording: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := StopRecordResponse{
		Message: "recording stopped for channel: " + req.Channel,
	}

	json.NewEncoder(w).Encode(response)
}
