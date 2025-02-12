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
	mux.Handle("/", http.FileServer(http.Dir("./front-end")))

	http.ListenAndServe(":8080", mux)
}
