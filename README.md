# Service Broker for a Cloud Foundry Marketplace

This project provides a Helm charge (a Kubernetes deployment of a single CLI, written in Golang within this repo) to run an HTTP-based API that implements the [Open Service Broker API](https://www.openservicebrokerapi.org/) to allow access to a service catalog available on an adjacent Cloud Foundry.

When you use Kubernetes you should not deploy, run, maintain, upgrade, backup, nor restore databases nor any stateful facilities. Leave it to people who will do it well. One way to separate these concerns is the Kubernetes incubator project [Service Catalog](https://svc-cat.io/). You gain access to a suite of "services" that your organization or third-party organizations are prepared to maintain for you as a service. For example, you might use the Service Catalog to request a PostgreSQL database from your underlying cloud provider.

This project provides a service broker to allow your Kubernetes cluster to access the pre-existing services registered with a neighboring Cloud Foundry marketplace. It only requires that your Kubernetes cluster has networking access to the Cloud Foundry API, and that your applications' pods have networking access to the provisioned service instances (such as databases).

## TODOs & Ideas

- [X] support UAA client/secret
- [x] support [GetInstance](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-instance)
- [x] support [GetBinding](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#fetching-a-service-binding) and update /v2/catalog
- [x] support async brokers with LastOperation
- [ ] support provision/bind parameters
- [ ] support Update (requires https://github.com/cloudfoundry-community/go-cfclient/issues/211)
- [ ] create target space if missing
- [ ] one space per kubernetes namespace
- [ ] accept named org/space and convert to GUIDs internally
- [ ] kubernetes/service catalog users mapped to backend Cloud Foundry users (perhaps with [Originating Identity](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#originating-identity))
- [ ] deprovision could also unbind all service keys
- [ ] app or pod to emit K8s events during start up

Edge cases:

- [ ] correctly consider [`?accepts_incomplete=false`](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#asynchronous-operations)
- [ ] [410 Gone on LastOperation](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md#polling-last-operation-for-service-instances)

Blocked:

- [ ] support async brokers with LastBindingOperation (blocked by https://github.com/cloudfoundry/cloud_controller_ng/issues/1246)

## Install/upgrade with Helm

You can configure the service broker to interact with a single Cloud Foundry API using a pre-existing user, or a UAA client (replace `cf.username` and `cf.password` values below with `cf.uaa_client_id` and `cf.uaa_client_secret` values).

Login to Cloud Foundry and create a space into which service instances will be created.

```shell
export CF_API=https://api.run.pivotal.io
export CF_USERNAME=...
export CF_PASSWORD=...
cf login -a $CF_API -u $CF_USERNAME -p $CF_PASSWORD

cf create-space playtime-cf-marketplace
cf target -s playtime-cf-marketplace
```

Next, install/upgrade the Helm chart:

```shell
helm plugin install https://github.com/hypnoglow/helm-s3.git
helm repo add starkandwayne s3://helm.starkandwayne.com/charts
helm repo update
helm upgrade --install pws-broker starkandwayne/cf-marketplace-servicebroker \
    --namespace catalog \
    --wait \
    --set "cf.api=$CF_API" \
    --set "cf.username=${CF_USERNAME:?required},cf.password=${CF_PASSWORD:?required}" \
    --set "cf.organizationGUID=$(jq -r .OrganizationFields.GUID ~/.cf/config.json)" \
    --set "cf.spaceGUID=$(jq -r .SpaceFields.GUID ~/.cf/config.json)"
```

Next, follow the instructions for registering with your Service Catalog. You'll now be able to view/provision/bind services within your Kubernetes cluster that are actually provisioned in the remote Cloud Foundry environment.

For example:

```shell
kubectl create secret generic pws-broker-cf-marketplace-servicebroker-basic-auth \
--from-literal username=broker \
--from-literal password=broker

svcat register pws-broker-cf-marketplace-servicebroker \
--url http://pws-broker-cf-marketplace-servicebroker.default.svc.cluster.local:8080 \
--scope cluster \
--basic-secret pws-broker-cf-marketplace-servicebroker-basic-auth
```

You'll now be able to view classes and plans, and to then instantiate and bind service instances.

```console
$ svcat get plans
               NAME                NAMESPACE                          CLASS                                   DESCRIPTION
+--------------------------------+-----------+-----------------------------------------------------+--------------------------------+
  trial                                        p-config-server                                       Service instances using this
                                                                                                     plan are deleted automatically
                                                                                                     7 days after creation
  standard                                     p-config-server                                       Standard Plan
  small                                        searchify                                             Small
  pro                                          searchify                                             Pro
  plus                                         searchify                                             Plus
  essential                                    amazon-s3-starkandwayne-optigit                       An S3 plan providing a single
                                                                                                     bucket with unlimited storage.
  standard                                     scheduler-for-pcf                                     Scheduler for PCF
```

## Dev/test

In one terminal, first configure for target Cloud Foundry and create a space into which service instances will be created:

```shell
export CF_API=https://api.run.pivotal.io
export CF_USERNAME=...
export CF_PASSWORD=...
cf login -a $CF_API -u $CF_USERNAME -p $CF_PASSWORD

cf create-space playtime-cf-marketplace
cf target -s playtime-cf-marketplace

export CF_ORGANIZATION_GUID=$(jq -r .OrganizationFields.GUID ~/.cf/config.json)
export CF_SPACE_GUID=$(jq -r .SpaceFields.GUID ~/.cf/config.json)
```

Next, run the broker.

From source:

```shell
go run cmd/cf-marketplace-servicebroker/main.go
```

From Docker image:

```sehll
docker run \
    -e CF_API=$CF_API \
    -e CF_USERNAME=$CF_USERNAME \
    -e CF_PASSWORD=$CF_PASSWORD \
    -p 8080:8080 \
    starkandwayne/cf-marketplace-servicebroker
```

In another terminal:

```example
curl -u broker:broker -H 'X-Broker-API-Version: 2.14' localhost:8080/v2/catalog
```

Or setup `eden`:

```bash
export SB_BROKER_URL=http://localhost:8080
export SB_BROKER_USERNAME=broker
export SB_BROKER_PASSWORD=broker
```

And see catalog:

```console
$ eden catalog
Service Name     Plan Name  Description
some-cf-service  plan-a     Probably smallest plan
```