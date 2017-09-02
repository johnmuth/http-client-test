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

	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.Info("About to get connection")
		},
		PutIdleConn: func(err error) {
			if err != nil {
				log.Info(fmt.Sprintf("Put idle connection failed. err=%v", err))
			} else {
				log.Info("Put idle connection succeeded.")
			}
		},
		Got100Continue : func() {
			log.Info("Got 100 Continue")
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Info("Got connection")
		},
		ConnectStart: func(network, addr string) {
			log.Info("Dial start")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.Info("DNS start", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			log.Info("DNS done")
		},
		ConnectDone: func(network, addr string, err error) {
			log.Info("Dial done")
		},
		GotFirstResponseByte: func() {
			log.Info("First response byte!")
		},
		WroteHeaders: func() {
			log.Info("Wrote headers")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			log.Info("Wrote request")
		},
	}
	var resp *http.Response
	log.Info(fmt.Sprintf("About to send request: POST %s with body %s", svc.BaseURL, serviceRequest.String()))
	req, err := http.NewRequest("POST", svc.BaseURL, strings.NewReader(serviceRequest.String()))
	if err != nil {
		log.Error("Error creating request to service", err)
		return
	}
	req.Header.Set("Content-type", "application/json")
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err = svc.HttpClient.Do(req)
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
