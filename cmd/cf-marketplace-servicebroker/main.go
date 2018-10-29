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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/broker"
	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/cfconfig"
	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/version"

	"code.cloudfoundry.org/lager"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/brokerapi"
)

func statusAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	servicebroker := broker.NewMarketplaceBrokerImpl()

	logger := lager.NewLogger("cf-marketplace-servicebroker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	cfconfig := cfconfig.NewConfigFromEnvVars()

	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	var netClient = &http.Client{
		Timeout: 3 * time.Second,
	}

	fmt.Printf("Connecting to Cloud Foundry %s...", cfconfig.API)
	cfclient, err := cf.NewClient(&cf.Config{
		ApiAddress:        cfconfig.API,
		Username:          cfconfig.Username,
		Password:          cfconfig.Password,
		Token:             cfconfig.AccessToken,
		SkipSslValidation: cfconfig.SSLSkipValidation,
		HttpClient:        netClient,
		UserAgent:         "cf-marketplace-servicebrokers/" + version.Version,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("OK!")

	fmt.Printf("\nFetching marketplace services...")
	cfServices, err := cfclient.ListServicesByQuery(url.Values{})
	if err != nil {
		panic(err)
	}
	fmt.Println("OK!")
	fmt.Printf("Found %d services\n", len(cfServices))

	fmt.Printf("Fetching service plans...")
	cfServicePlans, err := cfclient.ListServicePlans()
	if err != nil {
		panic(err)
	}
	fmt.Println("OK!")
	fmt.Printf("Found %d service plans\n", len(cfServicePlans))

	broker.Marketplace = make([]brokerapi.Service, len(cfServices))
	for i, cfService := range cfServices {
		broker.Marketplace[i].Name = cfService.Label
		broker.Marketplace[i].ID = cfService.Guid
		broker.Marketplace[i].Description = cfService.Description
		broker.Marketplace[i].Bindable = cfService.Bindable
		broker.Marketplace[i].Tags = cfService.Tags
		metadata := &brokerapi.ServiceMetadata{}
		err := json.Unmarshal([]byte(cfService.Extra), metadata)
		if err != nil {
			panic(err)
		}
		broker.Marketplace[i].Metadata = metadata
		plansCount := 0
		for _, cfPlan := range cfServicePlans {
			if cfService.Guid == cfPlan.ServiceGuid {
				plansCount++
			}
		}
		broker.Marketplace[i].Plans = make([]brokerapi.ServicePlan, plansCount)

		planIndex := 0
		for _, cfPlan := range cfServicePlans {
			if cfService.Guid == cfPlan.ServiceGuid {
				fmt.Printf("Adding plan %s to service %s\n", cfPlan.Name, cfService.Label)
				broker.Marketplace[i].Plans[planIndex] = brokerapi.ServicePlan{
					ID:          cfPlan.Guid,
					Name:        cfPlan.Name,
					Description: cfPlan.Description,
					Free:        &cfPlan.Free,
					Bindable:    &cfPlan.Bindable,
				}
				planIndex++
			}
		}
	}

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
