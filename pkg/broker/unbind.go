package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func (bkr *MarketplaceBrokerImpl) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (unbindSpec brokerapi.UnbindSpec, err error) {
	bkr.Logger.Info("unbind.start", lager.Data{"instanceID": instanceID, "bindID": bindingID})

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return unbindSpec, brokerapi.NewFailureResponse(err, 400, "cf-client")
	}

	cfSvcInstance, err := bkr.lookupServiceInstance(cfclient, instanceID)
	if err != nil {
		// Whilst it is an internal error, if the backend service doesn't exist anymore then allow Unbind to succeed.
		bkr.Logger.Error("lookup-service-instance", err, lager.Data{"instanceID": instanceID, "bindID": bindingID})
		return unbindSpec, nil
	}

	cfSvcKey, err := bkr.lookupServiceKey(cfclient, cfSvcInstance, bindingID)
	if err != nil {
		// Whilst it is an internal error, if the backend service key doesn't exist anymore then allow Unbind to succeed.
		bkr.Logger.Error("lookup-service-key", err, lager.Data{"instanceID": instanceID, "bindID": bindingID})
		return unbindSpec, nil
	}

	err = cfclient.DeleteServiceKey(cfSvcKey.Guid)
	if err != nil {
		return unbindSpec, brokerapi.NewFailureResponse(err, 400, "delete-service-key")
	}
	unbindSpec.IsAsync = false
	bkr.Logger.Info("unbind.end", lager.Data{"instanceID": instanceID, "bindID": bindingID})
	return
}
