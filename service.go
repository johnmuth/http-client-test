package main

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
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
	requestTimings := make(map[string]time.Time)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), clientTrace(serviceRequest, requestTimings)))
	log.WithField("requestid", serviceRequest.RequestID).Debug("About to send request to service")
	resp, err = svc.HttpClient.Do(req)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"requestTimings":  requestTimings,
			"requestid":  serviceRequest.RequestID,
			"error": err,
		}).Error("Error sending request to service")
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
	} else {
		log.WithFields(map[string]interface{}{
			"requestid": serviceRequest.RequestID,
		}).Debug("Done with request")
	}
	return
}

func clientTrace(serviceRequest ServiceRequest, requestTimings map[string]time.Time) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			requestTimings["getconn"] = time.Now()
			log.WithField("requestid", serviceRequest.RequestID).Debug("About to get connection")
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			requestTimings["gotconn"] = time.Now()
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"reused":    connInfo.Reused,
				"idletime":  connInfo.IdleTime,
				"wasidle":   connInfo.WasIdle,
			}).Debug("Got connection")
		},
		PutIdleConn: func(err error) {
			requestTimings["putidleconn"] = time.Now()
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"err":       err,
			}).Debug("Put idle connection")
		},
		Got100Continue: func() {
			requestTimings["got100continue"] = time.Now()
			log.WithField("requestid", serviceRequest.RequestID).Debug("Got 100 Continue")
		},
		ConnectStart: func(network, addr string) {
			requestTimings["connectstart"] = time.Now()
			log.WithField("requestid", serviceRequest.RequestID).Debug("Dial start")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			requestTimings["dnsstart"] = time.Now()
			log.WithField("requestid", serviceRequest.RequestID).Debug("DNS start", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			requestTimings["dnsdone"] = time.Now()
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"coalesced": info.Coalesced,
				"err":       info.Err,
			}).Debug("DNS done")
		},
		ConnectDone: func(network, addr string, err error) {
			requestTimings["connectdone"] = time.Now()
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"err":       err,
			}).Debug("Dial done")
		},
		GotFirstResponseByte: func() {
			requestTimings["gotfirstresponsebyte"] = time.Now()
			log.WithField("requestid", serviceRequest.RequestID).Debug("First response byte!")
		},
		WroteHeaders: func() {
			requestTimings["wroteheaders"] = time.Now()
			log.WithField("requestid", serviceRequest.RequestID).Debug("Wrote headers")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			requestTimings["wroterequest"] = time.Now()
			log.WithField("requestid", serviceRequest.RequestID).Debug("Wrote request")
		},
	}
}
