package cfconfig

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/version"

	cf "github.com/cloudfoundry-community/go-cfclient"
)

// Config contains the initial configuration for a target Cloud Foundry environment
type Config struct {
	API               string
	SSLSkipValidation bool
	Username          string
	Password          string
	UAAClientID       string
	UAAClientSecret   string
	AccessToken       string

	// Initially provision all services instances into one CF space
	OrganizationGUID string
	SpaceGUID        string
	BindingAppName   string

	HTTPClient *http.Client

	// Discovered/dynamically created state for the Cloud Foundry environment

	// Marketplace is the cache of services/plans offered by the target Cloud Foundry
	Marketplace []brokerapi.Service

	// AppGUID for the empty binding app (see cfconfig.CreateBindingApp)
	BindingAppGUID string
}

// NewConfigFromEnvVars constructs a Config from environment variables
func NewConfigFromEnvVars() (config *Config) {
	config = &Config{
		API:               os.Getenv("CF_API"),
		SSLSkipValidation: os.Getenv("CF_SKIP_SSL_VALIDATION") == "true",
		Username:          os.Getenv("CF_USERNAME"),
		Password:          os.Getenv("CF_PASSWORD"),
		UAAClientID:       os.Getenv("CF_UAA_CLIENT_ID"),
		UAAClientSecret:   os.Getenv("CF_UAA_CLIENT_SECRET"),
		AccessToken:       os.Getenv("CF_ACCESS_TOKEN"),

		OrganizationGUID: os.Getenv("CF_ORGANIZATION_GUID"),
		SpaceGUID:        os.Getenv("CF_SPACE_GUID"),

		BindingAppName: os.Getenv("CF_BINDING_APPNAME"),
	}
	if config.API == "" {
		fmt.Fprintln(os.Stderr, "ERROR: configure with $CF_API, and either $CF_ACCESS_TOKEN, or both $CF_USERNAME, $CF_PASSWORD")
		os.Exit(1)
	}
	if config.OrganizationGUID == "" || config.SpaceGUID == "" {
		fmt.Fprintln(os.Stderr, "ERROR: configure with both $CF_ORGANIZATION_GUID, and $CF_SPACE_GUID")
		os.Exit(1)
	}
	if config.BindingAppName == "" {
		config.BindingAppName = "cf-marketplace-servicebroker-binding-app"
	}
	return
}

func (config *Config) Client() (cfclient *cf.Client, err error) {
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 120 * time.Second,
		}
	}
	return cf.NewClient(&cf.Config{
		ApiAddress:        config.API,
		Username:          config.Username,
		Password:          config.Password,
		ClientID:          config.UAAClientID,
		ClientSecret:      config.UAAClientSecret,
		Token:             config.AccessToken,
		SkipSslValidation: config.SSLSkipValidation,
		HttpClient:        config.HTTPClient,
		UserAgent:         "cf-marketplace-servicebrokers/" + version.Version,
	})
}
