package main

type AppConfig struct {
	Port           int    `default:"8000"`
	ServiceBaseURL string `envconfig:"SERVICE_BASE_URL" required:"true"`
	Env            string `envconfig:"ENV_NAME" required:"true"`
}

func (c *AppConfig) IsLocal() bool {
	return c.Env == "local"
}
