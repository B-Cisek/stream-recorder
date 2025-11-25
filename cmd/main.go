package main

import (
	"log"
	"net/http"

	"github.com/B-Cisek/stream-recorder/internal/handlers"
	"github.com/B-Cisek/stream-recorder/internal/services"
)

func main() {
	recorderService := services.NewRecorderService()

	http.HandleFunc("/ping", handlers.NewPingHandler().Handle)
	http.HandleFunc("/record/start", handlers.NewStartRecordHandler(recorderService).Handle)
	http.HandleFunc("/record/stop", handlers.NewStopRecordHandler(recorderService).Handle)

	log.Println("Application started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
