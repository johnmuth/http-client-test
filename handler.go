package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/satori/go.uuid"
)

// Handler handles requests
type Handler struct {
	Service Service
}

// ServeHTTP serves HTTP
func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u1 := uuid.NewV4()
	serviceRequest := &ServiceRequest{UUID: u1.String()}
	log.WithField("uuid", serviceRequest.UUID).Info("About to do service.Call")
	serviceResponse, err := handler.Service.Call(*serviceRequest)
	log.WithField("uuid", serviceRequest.UUID).Info("Got response from service.Call")
	if err != nil {
		log.WithField("uuid", serviceRequest.UUID).Error("Error calling service", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, serviceResponse.String())
}
