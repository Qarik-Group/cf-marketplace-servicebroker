package broker

import (
	"context"
	"fmt"
	"net/url"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

// Deprovision forwards on a service instance deprovision request to the backend Cloud Foundry API
func (bkr *MarketplaceBrokerImpl) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (spec brokerapi.DeprovisionServiceSpec, err error) {
	bkr.Logger.Info("deprovision.start", lager.Data{"instanceID": instanceID})

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return
	}

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

	cfSvcInstance, err := bkr.lookupServiceInstance(cfclient, instanceID)
	if err != nil {
		return
	}

	err = cfclient.DeleteServiceInstance(cfSvcInstance.Guid, true, true)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	bkr.Logger.Info("deprovision.end", lager.Data{
		"instanceID": instanceID,
		"error":      errMsg,
	})
	return
}
