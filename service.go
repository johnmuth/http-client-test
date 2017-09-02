package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"errors"
	log "github.com/sirupsen/logrus"
)

// Service sends requests to a remote HTTP API
type Service struct {
	BaseURL    string
	HttpClient HttpClient
}

type HttpClient interface {
	Post(url string, bodyType string, body io.Reader) (*http.Response, error)
}

func (svc Service) Call(serviceRequest ServiceRequest) (serviceResponse ServiceResponse, err error) {
	serviceRequestURL := fmt.Sprintf("%s", svc.BaseURL)
	var resp *http.Response
	log.Info(fmt.Sprintf("About to send request: POST %s with body %s", serviceRequestURL, serviceRequest.String()))
	resp, err = svc.HttpClient.Post(serviceRequestURL, "application/json", strings.NewReader(serviceRequest.String()))
	if err != nil {
		log.Error("Error sending request to service", err)
		return
	}
	if resp.StatusCode != 200 {
		var respBody string
		if resp.Body != nil {
			respErrorBody, _ := ioutil.ReadAll(resp.Body)
			respBody = string(respErrorBody)
			resp.Body.Close()
		}
		errorMsg := fmt.Sprintf("Request to service returned status: %d and body: %s ", resp.StatusCode, respBody)
		log.Error(errorMsg)
		return serviceResponse, errors.New(errorMsg)
	}
	err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
	resp.Body.Close()
	if err != nil {
		log.Error("Error parsing response from Service.", err)
	}
	return
}
