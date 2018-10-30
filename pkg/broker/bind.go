package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

// Bind forwards on a service instance bind request to the backend Cloud Foundry API
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
	cfBindResp, err := cfclient.CreateServiceBinding(bkr.CF.BindingAppGUID, cfSvcInstance.Guid)
	if err != nil {
		return svcBinding, brokerapi.NewFailureResponse(err, 400, "create-binding")
	}
	svcBinding.Credentials = cfBindResp.Credentials

	bkr.Logger.Info("bind.end", lager.Data{
		"instanceID": instanceID,
		"bindID":     bindingID,
	})
	return
}
