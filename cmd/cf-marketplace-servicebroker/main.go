package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/broker"
	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/cfconfig"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func statusAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	logger := lager.NewLogger("cf-marketplace-servicebroker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	cf := cfconfig.NewConfigFromEnvVars()
	cf.DiscoverMarketplace()

	// TODO: The services that do not support service keys also don't support any binding
	// TODO: Test more services until we find at least 1 that supports binding but not service keys
	// TODO: else delete this code
	// cf.CreateBindingApp()

	servicebroker := broker.NewMarketplaceBrokerImpl(cf, logger)

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: "broker",
		Password: "broker",
	}
	brokerAPI := brokerapi.New(servicebroker, logger, brokerCredentials)

	http.HandleFunc("/health", statusAPI)
	http.Handle("/", brokerAPI)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("\n\nStarting Cloud Foundry Marketplace Broker on 0.0.0.0:" + port)
	logger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}
