package main

import (
	"log"
	"net/http"

	"github.com/B-Cisek/stream-recorder/internal/handlers"
	"github.com/B-Cisek/stream-recorder/internal/services"
)

func main() {
	recorderService := services.NewRecorderService()

	http.HandleFunc("/api/ping", handlers.NewPingHandler().Handle)
	http.HandleFunc("/api/record/start", handlers.NewStartRecordHandler(recorderService).Handle)
	http.HandleFunc("/api/record/stop", handlers.NewStopRecordHandler(recorderService).Handle)

	log.Println("Application started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
