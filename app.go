package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

func main() {

	config, err := LoadAppConfig()
	if err != nil {
		log.Error("Error loading config", err.Error())
	}

	log.Info("Listening on", config.Port)

	httpClient := &http.Client{
			// TODO: make these options configurable
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
			},
			Timeout: time.Millisecond * 1500,
		}

	service := &Service{config.ServiceBaseURL, httpClient}
	handler := &Handler{*service}
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
