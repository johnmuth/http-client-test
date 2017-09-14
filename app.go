package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"time"
	"github.com/viki-org/dnscache"
	"strings"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339Nano})

	config, err := LoadAppConfig()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error loading config")
	}

	log.WithField("port", config.Port).Info("Listening")

	resolver := dnscache.New(time.Second * 60) //how often to refresh cached dns records, happens in background

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: config.HTTPClientMaxIdleConnsPerHost,
			// Go does not cache DNS lookups, so we define a custom Dial function that does.
			// This fixed a problem where requests were timing out during DNS lookup
			// even though we were hitting the same hostname over and over.
			Dial: func(network string, address string) (net.Conn, error) {
				separator := strings.LastIndex(address, ":")
				ip, err := resolver.FetchOneString(address[:separator])
				if err != nil {
					return nil, err
				}
				return net.Dial("tcp", ip + address[separator:])
			},
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
		log.WithField("error", err.Error()).Error("Problem starting server")
		os.Exit(1)
	}
}

func LoadAppConfig() (*AppConfig, error) {
	var config AppConfig
	err := envconfig.Process("", &config)
	return &config, err
}
