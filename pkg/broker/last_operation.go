package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation looks up readiness/failure of asynchronous provision/update/deprovision operations
// https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#polling-last-operation-for-service-instances
func (bkr *MarketplaceBrokerImpl) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (spec brokerapi.LastOperation, err error) {
	bkr.Logger.Info("last-operation.start", lager.Data{"instanceID": instanceID})

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "cf-client")
	}

	cfSvcInstance, err := bkr.lookupServiceInstance(cfclient, instanceID)
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "lookup-service")
	}

	// https://apidocs.cloudfoundry.org/5.4.0/service_instances/retrieve_a_particular_service_instance.html
	// last_operation.state	- in progress, succeeded, or failed
	spec.State = brokerapi.LastOperationState(cfSvcInstance.LastOperation.State)
	spec.Description = cfSvcInstance.LastOperation.Description

	bkr.Logger.Info("last-operation.end", lager.Data{"instanceID": instanceID})
	return
}
