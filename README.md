# Service Broker for a Cloud Foundry Marketplace

## Install with Helm

Login to Cloud Foundry and create a space into which service instances will be created.

```shell
export CF_API=https://api.run.pivotal.io
export CF_USERNAME=...
export CF_PASSWORD=...
cf login -a $CF_API -u $CF_USERNAME -p $CF_PASSWORD

cf create-space playtime-cf-marketplace
cf target -s playtime-cf-marketplace
```

Next, config and install the Helm chart:

```shell
helm install ./helm --name pws-broker --wait \
    --set "cf.api=$CF_API" \
    --set "cf.username=${CF_USERNAME:?required},cf.password=${CF_PASSWORD:?required}" \
    --set "cf.organizationGUID=$(jq -r .OrganizationFields.GUID ~/.cf/config.json)" \
    --set "cf.spaceGUID=$(jq -r .SpaceFields.GUID ~/.cf/config.json)"
```

To upgrade, first login and target the space. Then run `helm upgrade`:

```shell
export CF_API=https://api.run.pivotal.io
export CF_USERNAME=...
export CF_PASSWORD=...
cf login -a $CF_API -u $CF_USERNAME -p $CF_PASSWORD
cf target -s playtime-cf-marketplace

helm upgrade pws-broker ./helm \
    --set "cf.api=$CF_API" \
    --set "cf.username=${CF_USERNAME:?required},cf.password=${CF_PASSWORD:?required}" \
    --set "cf.organizationGUID=$(jq -r .OrganizationFields.GUID ~/.cf/config.json)" \
    --set "cf.spaceGUID=$(jq -r .SpaceFields.GUID ~/.cf/config.json)"
```

Next, follow the instructions for registering with your Service Catalog. You'll now be able to view/provision/bind services within your Kubernetes cluster that are actually provisioned in the remote Cloud Foundry environment.

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
curl -u: -H 'X-Broker-API-Version: 2.12' localhost:8080/v2/catalog
```

Or setup `eden`:

```bash
export SB_BROKER_URL=http://localhost:8080
export SB_BROKER_USERNAME=
export SB_BROKER_PASSWORD=
```

And see catalog:

```console
$ eden catalog
Service Name     Plan Name  Description
some-cf-service  plan-a     Probably smallest plan
```