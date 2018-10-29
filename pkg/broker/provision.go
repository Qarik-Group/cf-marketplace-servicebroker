package broker

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"code.cloudfoundry.org/lager"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/brokerapi"
)

// Provision forwards on a service instance provision request to the backend Cloud Foundry API
func (bkr *MarketplaceBrokerImpl) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	bkr.Logger.Info("provision.start", lager.Data{"guid": instanceID})
	req := cf.ServiceInstanceRequest{
		Name:            instanceID,
		ServicePlanGuid: details.PlanID,
		SpaceGuid:       bkr.CF.SpaceGUID,
		// Parameters:      details.RawParameters,
	}

	cfclient, err := bkr.CF.Client()
	if err != nil {
		return
	}
	svcInstance, err := cfclient.CreateServiceInstance(req)
	if err != nil {
		return
	}

	spec.DashboardURL = svcInstance.DashboardUrl
	spec.IsAsync = svcInstance.LastOperation.State == "in progress"
	spec.OperationData = "provision::" + svcInstance.Guid

	bkr.Logger.Info("provision.end", lager.Data{
		"guid":                instanceID,
		"async":               spec.IsAsync,
		"operation-data":      spec.OperationData,
		"last-op.status":      svcInstance.LastOperation.State,
		"last-op.description": svcInstance.LastOperation.Description,
	})

	if svcInstance.LastOperation.State == "failed" {
		err = errors.New("cf failed to provision: " + svcInstance.LastOperation.Description)
	}

	return
}

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
