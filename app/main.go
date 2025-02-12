package main

import (
	"encoding/json"
	"net/http"
	"uk_realtime_bus/config"
	"uk_realtime_bus/services"
)

func main() {
	envConfig := config.NewEnvSettings()

	busLocatorService := services.NewBusLocatorService(envConfig)

	mux := http.NewServeMux()

	// TODO: create a handler service for such purposes
	mux.HandleFunc("/buses", func(w http.ResponseWriter, r *http.Request) {
		buses := busLocatorService.GetAllRealtimeBuses()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(buses)
	})
	// TODO fix that
	mux.Handle("/", http.FileServer(http.Dir("./app/static")))

	http.ListenAndServe(":8080", mux)
}
