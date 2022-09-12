#!/bin/bash

set -e -u

if [ -z "${PORT}" ]; then
  echo "Error: No PORT found" >&2
  exit 1
fi

# Database
db_credentials="$(echo "${VCAP_SERVICES}" | jq -r '.["mariadbent"][0].credentials // ""')"
if [ -z "${db_credentials}" ]; then
  echo "Error: Please bind a MariaDB service" >&2
  exit 1
fi
db_username="$(echo "${db_credentials}" | jq -r '.username // ""')"
db_password="$(echo "${db_credentials}" | jq -r '.password // ""')"
db_host="$(echo "${db_credentials}" | jq -r '.host // ""')"
db_port="$(echo "${db_credentials}" | jq -r '.port // ""')"
db_database="$(echo "${db_credentials}" | jq -r '.database // ""')"
db_string="${db_username}:${db_password}@(${db_host}:${db_port})/${db_database}"

# Run binary
./bin/clong \
    -port "${PORT}" \
    -db-string "${db_string}" \
    -force-https \
    -username "${USERNAME}" \
    -password "${PASSWORD}"
