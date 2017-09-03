package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// HTTPClientTestHandler handles requests
type HTTPClientTestHandler struct {
	Service Service
}

// ServeHTTP serves HTTP
func (handler HTTPClientTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u1 := uuid.NewV4()
	serviceRequest := &ServiceRequest{RequestID: u1.String()}
	log.WithField("requestid", serviceRequest.RequestID).Info("About to do service.Call")
	serviceResponse, err := handler.Service.Call(*serviceRequest)
	log.WithField("requestid", serviceRequest.RequestID).Info("Got response from service.Call")
	if err != nil {
		log.WithField("requestid", serviceRequest.RequestID).Error("Error calling service", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, serviceResponse.String())
}
