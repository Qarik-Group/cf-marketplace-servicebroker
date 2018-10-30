package broker

import (
	"context"

	"github.com/starkandwayne/cf-marketplace-servicebroker/pkg/cfconfig"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

// MarketplaceBrokerImpl describes the implementation of a broker of services registered to a single
// Cloud Foundry's marketplace
type MarketplaceBrokerImpl struct {
	CF     *cfconfig.Config
	Logger lager.Logger
}

// NewMarketplaceBrokerImpl creates a MarketplaceBrokerImpl
func NewMarketplaceBrokerImpl(cf *cfconfig.Config, logger lager.Logger) (bkr *MarketplaceBrokerImpl) {
	return &MarketplaceBrokerImpl{
		CF:     cf,
		Logger: logger,
	}
}

// Services creates the data returned by this Broker API's GET /v2/catalog endpoint
func (bkr *MarketplaceBrokerImpl) Services(ctx context.Context) (catalog []brokerapi.Service, err error) {
	return bkr.CF.Marketplace, nil
}

// LastOperation looks up readiness/failure of asynchronous operations
func (bkr *MarketplaceBrokerImpl) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
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
