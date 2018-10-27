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

	"github.com/starkandwayne/cf-marketplace-servicebroker/broker"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func main() {
	servicebroker := broker.NewMarketplaceBrokerImpl()

	logger := lager.NewLogger("cf-marketplace-servicebroker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

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
	fmt.Println("Starting Cloud Foundry Marketplace Broker on 0.0.0.0:" + port)
	logger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}
