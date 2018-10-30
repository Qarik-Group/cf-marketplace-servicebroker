package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

// GetBinding returns the service binding information again
// https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-binding
func (bkr *MarketplaceBrokerImpl) GetBinding(ctx context.Context, instanceID, bindingID string) (spec brokerapi.GetBindingSpec, err error) {
	bkr.Logger.Info("get-binding.start", lager.Data{"instanceID": instanceID, "bindingID": bindingID})

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "cf-client")
	}

	cfSvcInstance, err := bkr.lookupServiceInstance(cfclient, instanceID)
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "lookup-service")
	}

	svcKey, err := bkr.lookupServiceKey(cfclient, cfSvcInstance, bindingID)
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "lookup-service-key")
	}

	spec.Credentials = svcKey.Credentials
	// TODO: spec.Parameters = svcKey.

	bkr.Logger.Info("get-binding.end", lager.Data{"instanceID": instanceID, "bindingID": bindingID})
	return
}
