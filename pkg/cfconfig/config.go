package cfconfig

import (
	"fmt"
	"os"
)

// Config contains the initial configuration for a target Cloud Foundry environment
type Config struct {
	API               string
	SSLSkipValidation bool
	Username          string
	Password          string
	AccessToken       string
}

// NewConfigFromEnvVars constructs a Config from environment variables
func NewConfigFromEnvVars() (config *Config) {
	if os.Getenv("CF_API") == "" {
		fmt.Fprintln(os.Stderr, "ERROR: configure with $CF_API, and either $CF_ACCESS_TOKEN, or both $CF_USERNAME, $CF_PASSWORD")
		os.Exit(1)
	}
	return &Config{
		API:               os.Getenv("CF_API"),
		SSLSkipValidation: os.Getenv("CF_SKIP_SSL_VALIDATION") == "true",
		Username:          os.Getenv("CF_USERNAME"),
		Password:          os.Getenv("CF_PASSWORD"),
		AccessToken:       os.Getenv("CF_ACCESS_TOKEN"),
	}
}
