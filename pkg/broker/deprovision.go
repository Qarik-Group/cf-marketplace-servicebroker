package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

// Deprovision forwards on a service instance deprovision request to the backend Cloud Foundry API
func (bkr *MarketplaceBrokerImpl) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (spec brokerapi.DeprovisionServiceSpec, err error) {
	bkr.Logger.Info("deprovision.start", lager.Data{"instanceID": instanceID})

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "cf-client")
	}

	cfSvcInstance, err := bkr.lookupServiceInstance(cfclient, instanceID)
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "lookup-service")
	}

	err = cfclient.DeleteServiceInstance(cfSvcInstance.Guid, true, true)
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "delete-service")
	}

	bkr.Logger.Info("deprovision.end", lager.Data{"instanceID": instanceID})
	return
}
