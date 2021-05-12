#!/bin/sh
cd "$(dirname "$0")"

BIND_ADDRESS="${1-10.8.0.1}"

# remove if container already exists
docker stop "ioth-registry"
docker rm "ioth-registry"

# start new container registry
docker run -d \
    --name "ioth-registry" \
    --restart="always" \
    -p "${BIND_ADDRESS}:5000:5000" \
    -e "REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt" \
    -e "REGISTRY_HTTP_TLS_KEY=/certs/domain.key" \
    -v "${PWD}/certs:/certs" \
    registry:2

./apply.sh "${BIND_ADDRESS}"
