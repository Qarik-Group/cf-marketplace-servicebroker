package cfconfig

import (
	"fmt"
	"net/http"
	"os"

	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/version"

	cf "github.com/cloudfoundry-community/go-cfclient"
)

// Config contains the initial configuration for a target Cloud Foundry environment
type Config struct {
	API               string
	SSLSkipValidation bool
	Username          string
	Password          string
	AccessToken       string

	// Initially provision all services instances into one CF space
	OrganizationGUID string
	SpaceGUID        string

	HTTPClient *http.Client
}

// NewConfigFromEnvVars constructs a Config from environment variables
func NewConfigFromEnvVars() (config *Config) {
	config = &Config{
		API:               os.Getenv("CF_API"),
		SSLSkipValidation: os.Getenv("CF_SKIP_SSL_VALIDATION") == "true",
		Username:          os.Getenv("CF_USERNAME"),
		Password:          os.Getenv("CF_PASSWORD"),
		AccessToken:       os.Getenv("CF_ACCESS_TOKEN"),

		OrganizationGUID: os.Getenv("CF_ORGANIZATION_GUID"),
		SpaceGUID:        os.Getenv("CF_SPACE_GUID"),
	}
	if config.API == "" {
		fmt.Fprintln(os.Stderr, "ERROR: configure with $CF_API, and either $CF_ACCESS_TOKEN, or both $CF_USERNAME, $CF_PASSWORD")
		os.Exit(1)
	}
	if config.OrganizationGUID == "" || config.SpaceGUID == "" {
		fmt.Fprintln(os.Stderr, "ERROR: configure with both $CF_ORGANIZATION_GUID, and $CF_SPACE_GUID")
		os.Exit(1)
	}
	return
}

func (config *Config) Client() (cfclient *cf.Client, err error) {
	return cf.NewClient(&cf.Config{
		ApiAddress:        config.API,
		Username:          config.Username,
		Password:          config.Password,
		Token:             config.AccessToken,
		SkipSslValidation: config.SSLSkipValidation,
		HttpClient:        config.HTTPClient,
		UserAgent:         "cf-marketplace-servicebrokers/" + version.Version,
	})
}
