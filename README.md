# Service Broker for a Cloud Foundry Marketplace

## Install with Helm

```shell
export CF_API=https://api.run.pivotal.io
cf login -a $CF_API --sso

helm install ./helm --name pws-broker --wait \
    --set "cf.api=$CF_API,cf.accessToken=$(cf oauth-token | awk '{print $2}')"
```

To update/upgrade:

```shell
export CF_API=https://api.run.pivotal.io
cf login -a $CF_API --sso

helm upgrade pws-broker ./helm \
    --set "cf.api=$CF_API,cf.accessToken=$(cf oauth-token | awk '{print $2}')"
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

In one terminal, first configure for target Cloud Foundry:

```shell
export CF_API=https://api.run.pivotal.io
cf login -a $CF_API --sso

export CF_ACCESS_TOKEN="$(cf oauth-token | awk '{print $2}')"
```

Next, run the broker.

From source:

```shell
go run cmd/cf-marketplace-servicebroker/main.go
```

From Docker image:

```sehll
docker run \
    -e CF_API=https://api.run.pivotal.io \
    -e CF_ACCESS_TOKEN="$(cf oauth-token | awk '{print $2}')" \
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