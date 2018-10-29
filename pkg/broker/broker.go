package broker

import (
	"context"
	"errors"

	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/cfconfig"

	"code.cloudfoundry.org/lager"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/brokerapi"
)

// MarketplaceBrokerImpl describes the implementation of a broker of services registered to a single
// Cloud Foundry's marketplace
type MarketplaceBrokerImpl struct {
	CF     *cfconfig.Config
	Logger lager.Logger
}

// Marketplace is the cache of services/plans offered by the target Cloud Foundry
var Marketplace []brokerapi.Service

// NewMarketplaceBrokerImpl creates a MarketplaceBrokerImpl
func NewMarketplaceBrokerImpl(cf *cfconfig.Config, logger lager.Logger) (bkr *MarketplaceBrokerImpl) {
	return &MarketplaceBrokerImpl{
		CF:     cf,
		Logger: logger,
	}
}

// Services creates the data returned by this Broker API's GET /v2/catalog endpoint
func (bkr *MarketplaceBrokerImpl) Services(ctx context.Context) (catalog []brokerapi.Service, err error) {
	return Marketplace, nil
}

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
		"guid":                  instanceID,
		"async":                 spec.IsAsync,
		"operation-data":        spec.OperationData,
		"last-operation.status": svcInstance.LastOperation.State,
	})

	if svcInstance.LastOperation.State == "failed" {
		err = errors.New("cf failed to provision: " + svcInstance.LastOperation.Description)
	}

	return
}

// LastOperation looks up readiness/failure of asynchronous operations
func (bkr *MarketplaceBrokerImpl) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	panic("not implemented")
}

func (bkr *MarketplaceBrokerImpl) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	panic("not implemented")
}

func (bkr *MarketplaceBrokerImpl) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	panic("not implemented")
}

func (bkr *MarketplaceBrokerImpl) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	panic("not implemented")
}

func (bkr *MarketplaceBrokerImpl) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	panic("not implemented")
}

func (bkr *MarketplaceBrokerImpl) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	panic("not implemented")
}

func (bkr *MarketplaceBrokerImpl) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	panic("not implemented")
}

func (bkr *MarketplaceBrokerImpl) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	panic("not implemented")
}
