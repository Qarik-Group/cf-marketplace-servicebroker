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
	"net/url"
	"os"
	"time"

	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/broker"
	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/version"

	"code.cloudfoundry.org/lager"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/brokerapi"
)

func main() {
	servicebroker := broker.NewMarketplaceBrokerImpl()

	logger := lager.NewLogger("cf-marketplace-servicebroker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	if os.Getenv("CF_API") == "" {
		fmt.Fprintln(os.Stderr, "ERROR: configure with $CF_API, and either $CF_ACCESS_TOKEN, or both $CF_USERNAME, $CF_PASSWORD")
		os.Exit(1)
	}

	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	var netClient = &http.Client{
		Timeout: 3 * time.Second,
	}

	fmt.Printf("Connecting to Cloud Foundry %s...", os.Getenv("CF_API"))
	cfclient, err := cf.NewClient(&cf.Config{
		ApiAddress:        os.Getenv("CF_API"),
		Username:          os.Getenv("CF_USERNANE"),
		Password:          os.Getenv("CF_PASSWORD"),
		Token:             os.Getenv("CF_ACCESS_TOKEN"),
		SkipSslValidation: os.Getenv("CF_SKIP_SSL_VALIDATION") == "true",
		HttpClient:        netClient,
		UserAgent:         "cf-marketplace-servicebrokers/" + version.Version,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("done!")

	fmt.Printf("\nFetching marketplace services...")
	cfServices, err := cfclient.ListServicesByQuery(url.Values{})
	if err != nil {
		panic(err)
	}
	fmt.Println("done!")

	fmt.Printf("Found %d services\n", len(cfServices))

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: "",
		Password: "",
	}

	brokerAPI := brokerapi.New(servicebroker, logger, brokerCredentials)

	http.Handle("/", brokerAPI)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("\n\nStarting Cloud Foundry Marketplace Broker on 0.0.0.0:" + port)
	logger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}
