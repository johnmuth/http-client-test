package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Handler handles requests
type Handler struct {
	Service Service
}

// ServeHTTP serves HTTP
func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serviceRequest := &ServiceRequest{Foo:"abc", Bar:"def"}
	serviceResponse, err := handler.Service.Call(*serviceRequest)
	if err != nil {
		log.Error("Error calling service", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, serviceResponse.String())
}
