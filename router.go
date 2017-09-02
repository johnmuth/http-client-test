package main

import (
	"net/http"
	"fmt"
)

func NewRouter(handler *Handler) http.Handler {

	serveMux := http.NewServeMux()

	// Add healthcheck handler

	serveMux.HandleFunc("/internal/healthcheck", InternalHealthCheck)

	// Add api handler
	apiPath := "/api"
	serveMux.Handle(apiPath, handler)
	return serveMux

}

func InternalHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Healthy")
}
