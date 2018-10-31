package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/brokerapi"
)

// Bind forwards on a service instance bind request to the backend Cloud Foundry API
// Bind is implemented with CF Service Keys. Currently these do not support asynchronous bindings https://github.com/cloudfoundry/cloud_controller_ng/issues/1246
func (bkr *MarketplaceBrokerImpl) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (svcBinding brokerapi.Binding, err error) {
	bkr.Logger.Info("bind.start", lager.Data{"instanceID": instanceID, "bindID": bindingID})

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return svcBinding, brokerapi.NewFailureResponse(err, 400, "cf-client")
	}

	cfSvcInstance, err := bkr.lookupServiceInstance(cfclient, instanceID)
	if err != nil {
		return svcBinding, brokerapi.NewFailureResponse(err, 400, "lookup-service")
	}

	// TODO: The services that do not support service keys also don't support any binding
	// TODO: Test more services until we find at least 1 that supports binding but not service keys
	// TODO: else delete this code
	// cfBindResp, err := cfclient.CreateServiceBinding(bkr.CF.BindingAppGUID, cfSvcInstance.Guid)
	// if err != nil {
	// 	return svcBinding, brokerapi.NewFailureResponse(err, 400, "create-binding")
	// }
	// svcBinding.Credentials = cfBindResp.Credentials

	svcKeyReq := cf.CreateServiceKeyRequest{
		ServiceInstanceGuid: cfSvcInstance.Guid,
		Name:                bindingID,
		// Parameters:
	}
	svcKey, err := cfclient.CreateServiceKey(svcKeyReq)
	if err != nil {
		return svcBinding, brokerapi.NewFailureResponse(err, 400, "create-service-key")
	}

	svcBinding.Credentials = svcKey.Credentials

	bkr.Logger.Info("bind.end", lager.Data{"instanceID": instanceID, "bindID": bindingID})
	return
}
