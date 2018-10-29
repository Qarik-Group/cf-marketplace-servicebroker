package broker

import (
	"context"

	"github.com/pivotal-cf/brokerapi"
)

// MarketplaceBrokerImpl describes the implementation of a broker of services registered to a single
// Cloud Foundry's marketplace
type MarketplaceBrokerImpl struct {
}

// Marketplace is the cache of services/plans offered by the target Cloud Foundry
var Marketplace []brokerapi.Service

// NewMarketplaceBrokerImpl creates a MarketplaceBrokerImpl
func NewMarketplaceBrokerImpl() (bkr *MarketplaceBrokerImpl) {
	return &MarketplaceBrokerImpl{}
}

// Services creates the data returned by Broker API's GET /v2/catalog endpoint
func (bkr *MarketplaceBrokerImpl) Services(ctx context.Context) (catalog []brokerapi.Service, err error) {
	return Marketplace, nil
}

func (bkr *MarketplaceBrokerImpl) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
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

func (bkr *MarketplaceBrokerImpl) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	panic("not implemented")
}

func (bkr *MarketplaceBrokerImpl) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	panic("not implemented")
}
