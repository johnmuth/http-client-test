package main

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strings"
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
	serviceResponse.RequestID = serviceRequest.RequestID
	var resp *http.Response
	req, err := http.NewRequest("POST", svc.BaseURL, strings.NewReader(serviceRequest.String()))
	if err != nil {
		log.WithField("requestid", serviceRequest.RequestID).Error("Error creating request to service", err)
		return
	}
	req.Header.Set("Content-type", "application/json")
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), clientTrace(serviceRequest)))
	log.WithField("requestid", serviceRequest.RequestID).Info("About to send request to service")
	resp, err = svc.HttpClient.Do(req)
	if err != nil {
		log.WithField("requestid", serviceRequest.RequestID).Error("Error sending request to service", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var respBody string
		if resp.Body != nil {
			respErrorBody, _ := ioutil.ReadAll(resp.Body)
			respBody = string(respErrorBody)
		}
		log.WithFields(map[string]interface{}{
			"statuscode": resp.StatusCode,
			"body":       respBody,
			"requestid":  serviceRequest.RequestID,
		}).Error("Service returned non-200 response")
		return serviceResponse, errors.New("Service returned non-200 response")
	}
	err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"err":       err,
			"requestid": serviceRequest.RequestID,
		}).Error("Error parsing response from Service.")
	}
	return
}

func clientTrace(serviceRequest ServiceRequest) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.WithField("requestid", serviceRequest.RequestID).Info("About to get connection")
		},
		PutIdleConn: func(err error) {
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"err":       err,
			}).Info("Put idle connection")
		},
		Got100Continue: func() {
			log.WithField("requestid", serviceRequest.RequestID).Info("Got 100 Continue")
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"reused":    connInfo.Reused,
				"idletime":  connInfo.IdleTime,
				"wasidle":   connInfo.WasIdle,
			}).Info("Got connection")
		},
		ConnectStart: func(network, addr string) {
			log.WithField("requestid", serviceRequest.RequestID).Info("Dial start")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.WithField("requestid", serviceRequest.RequestID).Info("DNS start", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"coalesced": info.Coalesced,
				"err":       info.Err,
			}).Info("DNS done")
		},
		ConnectDone: func(network, addr string, err error) {
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"err":       err,
			}).Info("Dial done")
		},
		GotFirstResponseByte: func() {
			log.WithField("requestid", serviceRequest.RequestID).Info("First response byte!")
		},
		WroteHeaders: func() {
			log.WithField("requestid", serviceRequest.RequestID).Info("Wrote headers")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			log.WithField("requestid", serviceRequest.RequestID).Info("Wrote request")
		},
	}
}
