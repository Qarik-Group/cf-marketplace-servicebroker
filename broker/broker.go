package broker

import (
	"context"

	"github.com/pivotal-cf/brokerapi"
)

// MarketplaceBrokerImpl describes the implementation of a broker of services registered to a single
// Cloud Foundry's marketplace
type MarketplaceBrokerImpl struct {
}

func NewMarketplaceBrokerImpl() (bkr *MarketplaceBrokerImpl) {
	return &MarketplaceBrokerImpl{}
}

func (bkr *MarketplaceBrokerImpl) Services(ctx context.Context) (catalog []brokerapi.Service, err error) {
	catalog = []brokerapi.Service{
		brokerapi.Service{
			Name:        "some-cf-service",
			Description: "This is a CF service we found for you to use",
			ID:          "4f618a0d-8f29-42e2-a3b1-9ff294726226",
			Bindable:    true,
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          "bb0c42c8-66ad-4d82-929d-172dc415fea8",
					Name:        "plan-a",
					Description: "Probably smallest plan",
				},
			},
		},
	}
	return
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
