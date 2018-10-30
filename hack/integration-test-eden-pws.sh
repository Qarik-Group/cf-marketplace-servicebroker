#!/bin/bash

set -eu

[[ "${SB_BROKER_URL:-X}" == "X" ]] && { >&2 echo "ERROR: set \$SB_BROKER_URL to broker"; exit 1; }
export SB_BROKER_USERNAME=${SB_BROKER_USERNAME:-broker}
export SB_BROKER_PASSWORD=${SB_BROKER_PASSWORD:-broker}
export EDEN_CONFIG=$(mktemp -d)/eden.config

cf_api=$(bosh int ~/.cf/config.json --path /Target)
cf_space_name=$(bosh int ~/.cf/config.json --path /SpaceFields/Name)
[[ "${cf_api}" == "https://api.run.pivotal.io" ]] || { >&2 echo "ERROR: please 'cf login -a 'https://api.run.pivotal.io' first"; exit 1; }
[[ "${cf_space_name}" == "playtime-cf-marketplace" ]] || { >&2 echo "ERROR: please 'cf target -s playtime-cf-marketplace' first"; exit 1; }

eden catalog

echo
echo "** Provisioning elephantsql/turtle..."
eden provision -s elephantsql -p turtle

instanceID=$(bosh int $EDEN_CONFIG --path /service_instances/0/id)
function finish {
  set +e
  echo
  echo "** Review of Cloud Foundry assets remaining..."
  cf services
}
trap finish EXIT TERM QUIT INT

function curlBroker {
  path=$1; shift
  curl -sSf -H 'X-Broker-API-Version: 2.14' -u "${SB_BROKER_USERNAME}:${SB_BROKER_PASSWORD}" ${SB_BROKER_URL}${path} "$@"
}

echo
echo "** API GetInstance..."
curlBroker /v2/service_instances/${instanceID}

echo
echo "** Create binding for elephantsql/turtle..."
eden bind -i ${instanceID}

echo
echo "** View Cloud Foundry service key created"
cf service-keys ${instanceID}

echo
echo "** View eden config"
cat $EDEN_CONFIG
bindingID=$(bosh int $EDEN_CONFIG --path /service_instances/0/bindings/0/id)

echo
echo "** API GetBinding..."
curlBroker /v2/service_instances/${instanceID}/service_bindings/${bindingID}

echo
echo "** Unbinding from elephantsql/turtle..."
eden unbind -i ${instanceID} -b ${bindingID}

echo
echo "** Deprovisioning elephantsql/turtle (${instanceID})..."
eden deprovision -i ${instanceID}
