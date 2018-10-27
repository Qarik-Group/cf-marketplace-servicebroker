# Service Broker for a Cloud Foundry Marketplace

## Dev/test

In one terminal:

```console
$ go run cmd/cf-marketplace-servicebroker/main.go
Starting Cloud Foundry Marketplace Broker...
```

In another:

```console
$ curl -u: -H 'X-Broker-API-Version: 2.12' localhost:8080/v2/catalog
{"services":[]}
```

Or setup `eden`:

```bash
export SB_BROKER_URL=http://localhost:8080
export SB_BROKER_USERNAME=
export SB_BROKER_PASSWORD=
```

And see empty catalog:

```console
$ eden catalog
Service Name  Plan Name  Description

0 services
```