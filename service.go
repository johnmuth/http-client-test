package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http/httptrace"
)

// Service sends requests to a remote HTTP API
type Service struct {
	BaseURL    string
	HttpClient HttpClient
}

type HttpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

func (svc Service) Call(serviceRequest ServiceRequest) (serviceResponse ServiceResponse, err error) {
	var resp *http.Response
	req, err := http.NewRequest("POST", svc.BaseURL, strings.NewReader(serviceRequest.String()))
	if err != nil {
		log.WithField("uuid", serviceRequest.UUID).Error("Error creating request to service", err)
		return
	}
	req.Header.Set("Content-type", "application/json")
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), clientTrace(serviceRequest)))
	log.WithField("uuid", serviceRequest.UUID).Info("About to send request to service")
	resp, err = svc.HttpClient.Do(req)
	if err != nil {
		log.WithField("uuid", serviceRequest.UUID).Error("Error sending request to service", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var respBody string
		if resp.Body != nil {
			respErrorBody, _ := ioutil.ReadAll(resp.Body)
			respBody = string(respErrorBody)
		}
		errorMsg := fmt.Sprintf("Request to service returned status: %d and body: %s ", resp.StatusCode, respBody)
		log.WithField("uuid", serviceRequest.UUID).Error(errorMsg)
		return serviceResponse, errors.New(errorMsg)
	}
	err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
	if err != nil {
		log.WithField("uuid", serviceRequest.UUID).Error("Error parsing response from Service.", err)
	}
	serviceResponse.UUID = serviceRequest.UUID
	return
}

func clientTrace(serviceRequest ServiceRequest) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.WithField("uuid", serviceRequest.UUID).Info("About to get connection")
		},
		PutIdleConn: func(err error) {
			log.WithFields(map[string]interface{}{
				"uuid": serviceRequest.UUID,
				"err": err,
			}).Info("Put idle connection")
		},
		Got100Continue : func() {
			log.WithField("uuid", serviceRequest.UUID).Info("Got 100 Continue")
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.WithFields(map[string]interface{}{
				"uuid": serviceRequest.UUID,
				"reused": connInfo.Reused,
				"idletime": connInfo.IdleTime,
				"wasidle": connInfo.WasIdle,
			}).Info("Got connection")
		},
		ConnectStart: func(network, addr string) {
			log.WithField("uuid", serviceRequest.UUID).Info("Dial start")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.WithField("uuid", serviceRequest.UUID).Info("DNS start", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			log.WithFields(map[string]interface{}{
				"uuid": serviceRequest.UUID,
				"coalesced": info.Coalesced,
				"err": info.Err,
			}).Info("DNS done")
		},
		ConnectDone: func(network, addr string, err error) {
			log.WithFields(map[string]interface{}{
				"uuid": serviceRequest.UUID,
				"err": err,
			}).Info("Dial done")
		},
		GotFirstResponseByte: func() {
			log.WithField("uuid", serviceRequest.UUID).Info("First response byte!")
		},
		WroteHeaders: func() {
			log.WithField("uuid", serviceRequest.UUID).Info("Wrote headers")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			log.WithField("uuid", serviceRequest.UUID).Info("Wrote request")
		},
	}
}