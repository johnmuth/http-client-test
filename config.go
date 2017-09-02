package main

type AppConfig struct {
	Port                          int    `default:"8000"`
	ServiceBaseURL                string `envconfig:"SERVICE_BASE_URL" required:"true"`
	Env                           string `envconfig:"ENV_NAME" required:"true"`
	HTTPClientMaxIdleConnsPerHost int `envconfig:"HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST" required:"true"`
	HTTPClientMaxIdleConns int `envconfig:"HTTP_CLIENT_MAX_IDLE_CONNS" required:"true"`
	HTTPClientDialerTimeoutMS int `envconfig:"HTTP_CLIENT_DIALER_TIMEOUT_MS" required:"true"`
	HTTPClientDialerKeepAliveMS int `envconfig:"HTTP_CLIENT_DIALER_KEEPALIVE_MS" required:"true"`
	HTTPClientIdleConnTimeoutMS int `envconfig:"HTTP_CLIENT_IDLE_CONN_TIMEOUT_MS" required:"true"`
	HTTPClientTLSHandshakeTimeoutMS int `envconfig:"HTTP_CLIENT_TLS_HANDSHAKE_TIMEOUT_MS" required:"true"`
	HTTPClientExpectContinueTimeoutMS int `envconfig:"HTTP_CLIENT_EXPECT_CONTINUE_TIMEOUT_MS" required:"true"`
	HTTPClientTimeoutMS int `envconfig:"HTTP_CLIENT_TIMEOUT_MS" required:"true"`
}

func (c *AppConfig) IsLocal() bool {
	return c.Env == "local"
}
