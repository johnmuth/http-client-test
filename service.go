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

type requestTiming struct {
	time time.Time
	phase string
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
	requestTimings := []requestTiming{}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), clientTrace(serviceRequest, requestTimings)))
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
	} else {
		log.WithFields(map[string]interface{}{
			"requestid": serviceRequest.RequestID,
		}).Info("Done with request")
	}
	return
}

func clientTrace(serviceRequest ServiceRequest, requestTimings []requestTiming) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"GetConn"})
			log.WithField("requestid", serviceRequest.RequestID).Info("About to get connection")
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"GotConn"})
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"reused":    connInfo.Reused,
				"idletime":  connInfo.IdleTime,
				"wasidle":   connInfo.WasIdle,
			}).Info("Got connection")
		},
		PutIdleConn: func(err error) {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"PutIdleConn"})
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"err":       err,
			}).Info("Put idle connection")
		},
		Got100Continue: func() {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"Got100Continue"})
			log.WithField("requestid", serviceRequest.RequestID).Info("Got 100 Continue")
		},
		ConnectStart: func(network, addr string) {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"ConnectStart"})
			log.WithField("requestid", serviceRequest.RequestID).Info("Dial start")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"DNSStart"})
			log.WithField("requestid", serviceRequest.RequestID).Info("DNS start", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"DNSDone"})
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"coalesced": info.Coalesced,
				"err":       info.Err,
			}).Info("DNS done")
		},
		ConnectDone: func(network, addr string, err error) {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"ConnectDone"})
			log.WithFields(map[string]interface{}{
				"requestid": serviceRequest.RequestID,
				"err":       err,
			}).Info("Dial done")
		},
		GotFirstResponseByte: func() {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"GotFirstResponseByte"})
			log.WithField("requestid", serviceRequest.RequestID).Info("First response byte!")
		},
		WroteHeaders: func() {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"WroteHeaders"})
			log.WithField("requestid", serviceRequest.RequestID).Info("Wrote headers")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			requestTimings = append(requestTimings, requestTiming{time.Now(),"WroteRequest"})
			log.WithField("requestid", serviceRequest.RequestID).Info("Wrote request")
		},
	}
}
