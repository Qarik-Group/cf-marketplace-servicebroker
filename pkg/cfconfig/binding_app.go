package cfconfig

import (
	"fmt"
	"net/url"

	cf "github.com/cloudfoundry-community/go-cfclient"
)

// Some backend Cloud Foundry services might not support service keys.
// To ensure that we can generate binding credentials for any arbitrary service
// we will create a dummy application into the Cloud Foundry space (without pushing code)
// and use normal service bindings; rather than optionally supported service keys.

// CreateBindingApp creates a dummy Cloud Foundry application within the target space.
// If application already exists, it does nothing.
func (config *Config) CreateBindingApp() {
	fmt.Printf("Detecting if empty binding application already exists in space...")
	cfclient, err := config.Client()
	if err != nil {
		panic(err)
	}

	// 1. Look up config.BindingAppName by name in config.SpaceGUID
	app, err := config.findAppByName(cfclient, config.BindingAppName, config.SpaceGUID)
	if err != nil {
		panic(err)
	}
	if app != nil {
		fmt.Println("FOUND!")
		config.BindingAppGUID = app.Guid
		return
	}
	fmt.Println("CREATING!")

	// 2. If not found, create app
	app, err = config.createEmptyBindingApp(cfclient, config.BindingAppName, config.SpaceGUID)
	if err != nil {
		panic(err)
	}
	config.BindingAppGUID = app.Guid
}

// lookupServiceInstance converts an incoming service instance ID into the
// GUID for a backend Cloud Foundry service instance
//
// The incoming service instance ID was originally used in Provision(instanceID)
// as the name of the Cloud Foundry service instance.
// This method looks up the service instance by name, where instanceID is the name.
func (config *Config) findAppByName(cfclient *cf.Client, appName, spaceGUID string) (app *cf.App, err error) {
	query := url.Values{
		"q": []string{
			fmt.Sprintf("name:%s", appName),
			fmt.Sprintf("space_guid:%s", config.SpaceGUID),
		},
	}
	apps, err := cfclient.ListAppsByQuery(query)
	if err == nil && len(apps) == 0 {
		return nil, nil
	}
	if err == nil && len(apps) > 1 {
		err = fmt.Errorf("Too many Cloud Foundry applications found for query '%s'", query.Encode())
	}
	if err != nil {
		return
	}

	return &apps[0], nil
}

// createEmptyBindingApp creates an application without pushing or running any code
func (config *Config) createEmptyBindingApp(cfclient *cf.Client, appName, spaceGUID string) (*cf.App, error) {
	req := cf.AppCreateRequest{
		Name:        appName,
		SpaceGuid:   spaceGUID,
		DockerImage: "nginx",
	}
	app, err := cfclient.CreateApp(req)
	return &app, err
}
