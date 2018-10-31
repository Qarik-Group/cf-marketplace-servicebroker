package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

// GetInstance returns the service instance information again
// https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-instance
func (bkr *MarketplaceBrokerImpl) GetInstance(ctx context.Context, instanceID string) (spec brokerapi.GetInstanceDetailsSpec, err error) {
	bkr.Logger.Info("get-instance.start", lager.Data{"instanceID": instanceID})

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "cf-client")
	}

	cfSvcInstance, err := bkr.lookupServiceInstance(cfclient, instanceID)
	if err != nil {
		return spec, brokerapi.NewFailureResponse(err, 400, "lookup-service")
	}

	spec.PlanID = cfSvcInstance.ServicePlanGuid
	spec.ServiceID = cfSvcInstance.ServiceGuid
	spec.DashboardURL = cfSvcInstance.DashboardUrl
	// TODO: spec.Parameters = cfSvcInstance.
	// https://apidocs.cloudfoundry.org/5.4.0/service_instances/retrieve_a_particular_service_instance_parameters_experimental.html

	bkr.Logger.Info("get-instance.end", lager.Data{"instanceID": instanceID})
	return
}
