/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	fmt.Printf("Connecting to Cloud Foundry %s...", cf.API)

	cf.DiscoverCloudFoundryMarketplace()

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: "broker",
		Password: "broker",
	}

	servicebroker := broker.NewMarketplaceBrokerImpl(cf, logger)

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
