package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/brokerapi"
)

// Bind forwards on a service instance bind request to the backend Cloud Foundry API
func (bkr *MarketplaceBrokerImpl) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (svcBinding brokerapi.Binding, err error) {
	bkr.Logger.Info("bind.start", lager.Data{"instanceID": instanceID, "bindID": bindingID})

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return
	}

	cfSvcInstance, err := bkr.lookupServiceInstance(cfclient, instanceID)
	if err != nil {
		return
	}

	svcKeyReq := cf.CreateServiceKeyRequest{
		ServiceInstanceGuid: cfSvcInstance.Guid,
		Name:                bindingID,
		// Parameters: details.GetRawParameters(),
	}
	var svcKey cf.ServiceKey
	svcKey, err = cfclient.CreateServiceKey(svcKeyReq)
	svcBinding.Credentials = svcKey.Credentials

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	bkr.Logger.Info("deprovision.end", lager.Data{
		"instanceID": instanceID,
		"bindID":     bindingID,
		"error":      errMsg,
	})
	return
}
