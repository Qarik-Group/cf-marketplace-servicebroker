# Service Broker for a Cloud Foundry Marketplace

## Dev/test

In one terminal, first configure for target Cloud Foundry:

```commands
export CF_API=https://api.run.pivotal.io
cf login -a $CF_API --sso

export CF_ACCESS_TOKEN="$(cf oauth-token | awk '{print $2}')"
```

Next, run the broker:

```console
$ go run cmd/cf-marketplace-servicebroker/main.go
Starting Cloud Foundry Marketplace Broker...
```

In another:

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