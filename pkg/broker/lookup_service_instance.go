package broker

import (
	"fmt"
	"net/url"

	cf "github.com/cloudfoundry-community/go-cfclient"
)

// lookupServiceInstance converts an incoming service instance ID into the
// GUID for a backend Cloud Foundry service instance
//
// The incoming service instance ID was originally used in Provision(instanceID)
// as the name of the Cloud Foundry service instance.
// This method looks up the service instance by name, where instanceID is the name.
func (bkr *MarketplaceBrokerImpl) lookupServiceInstance(cfclient *cf.Client, instanceID string) (svcInstance *cf.ServiceInstance, err error) {
	// the incoming instanceID is actually a CF service instance name,
	// so first need to convert "instanceID" into CF service instance GUID
	query := url.Values{
		"q": []string{
			fmt.Sprintf("name:%s", instanceID),
			fmt.Sprintf("space_guid:%s", bkr.CF.SpaceGUID),
		},
	}
	svcInstances, err := cfclient.ListServiceInstancesByQuery(query)
	if err == nil && len(svcInstances) == 0 {
		err = fmt.Errorf("No Cloud Foundry service instance found for ID %s", instanceID)
	}
	if err == nil && len(svcInstances) > 1 {
		err = fmt.Errorf("Too many Cloud Foundry service instances found for query '%s'", query.Encode())
	}
	if err != nil {
		return
	}

	svcInstance = &svcInstances[0]
	return
}
