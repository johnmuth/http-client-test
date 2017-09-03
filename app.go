package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339Nano})

	config, err := LoadAppConfig()
	if err != nil {
		log.Error("Error loading config", err.Error())
	}

	log.Info("Listening on", config.Port)

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: config.HTTPClientMaxIdleConnsPerHost,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(config.HTTPClientDialerTimeoutMS) * time.Millisecond,
				KeepAlive: time.Duration(config.HTTPClientDialerKeepAliveMS) * time.Millisecond,
			}).DialContext,
			MaxIdleConns:          config.HTTPClientMaxIdleConns,
			IdleConnTimeout:       time.Duration(config.HTTPClientIdleConnTimeoutMS) * time.Millisecond,
			TLSHandshakeTimeout:   time.Duration(config.HTTPClientTLSHandshakeTimeoutMS) * time.Millisecond,
			ExpectContinueTimeout: time.Duration(config.HTTPClientExpectContinueTimeoutMS) * time.Millisecond,
		},
		Timeout: time.Duration(config.HTTPClientTimeoutMS) * time.Millisecond,
	}

	service := &Service{config.ServiceBaseURL, httpClient}
	handler := &HTTPClientTestHandler{*service}
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), NewRouter(handler))

	if err != nil {
		log.Error("Problem starting server", err.Error())
		os.Exit(1)
	}
}

func LoadAppConfig() (*AppConfig, error) {
	var config AppConfig
	err := envconfig.Process("", &config)
	return &config, err
}
