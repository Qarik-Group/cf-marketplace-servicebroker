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
	bkr.Logger.Info("deprovision.start", lager.Data{"guid": instanceID})

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

	cfSvcInstanceID := svcInstances[0].Guid

	err = cfclient.DeleteServiceInstance(cfSvcInstanceID, true, true)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	bkr.Logger.Info("deprovision.end", lager.Data{
		"guid":  instanceID,
		"error": errMsg,
	})
	return
}
