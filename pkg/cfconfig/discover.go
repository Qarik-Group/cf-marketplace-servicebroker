package cfconfig

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pivotal-cf/brokerapi"
)

// DiscoverMarketplace fetches all Services & Plans and
// constructs this OSBAPI /v2/catalog of provided services
//
// When written, it was assumed this function was only run once during start up.
// So it exits the application rather than return an error.
func (config *Config) DiscoverMarketplace() {
	fmt.Printf("Connecting to Cloud Foundry %s...", config.API)
	cfclient, err := config.Client()
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

	config.Marketplace = make([]brokerapi.Service, len(cfServices))
	for i, cfService := range cfServices {
		config.Marketplace[i].Name = cfService.Label
		config.Marketplace[i].ID = cfService.Guid
		config.Marketplace[i].Description = cfService.Description
		config.Marketplace[i].Bindable = cfService.Bindable
		config.Marketplace[i].Tags = cfService.Tags
		config.Marketplace[i].InstancesRetrievable = true
		config.Marketplace[i].BindingsRetrievable = true
		metadata := &brokerapi.ServiceMetadata{}
		err := json.Unmarshal([]byte(cfService.Extra), metadata)
		if err != nil {
			panic(err)
		}
		config.Marketplace[i].Metadata = metadata
		plansCount := 0
		for _, cfPlan := range cfServicePlans {
			if cfService.Guid == cfPlan.ServiceGuid {
				plansCount++
			}
		}
		config.Marketplace[i].Plans = make([]brokerapi.ServicePlan, plansCount)

		planIndex := 0
		for _, cfPlan := range cfServicePlans {
			if cfService.Guid == cfPlan.ServiceGuid {
				fmt.Printf("Adding plan %s to service %s\n", cfPlan.Name, cfService.Label)
				config.Marketplace[i].Plans[planIndex] = brokerapi.ServicePlan{
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
}
