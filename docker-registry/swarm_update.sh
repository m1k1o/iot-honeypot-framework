#!/bin/sh
cd "$(dirname "$0")"

BIND_ADDRESS="${1}"

docker build -t "ioth-registry:5000/registry-update-certs:latest" .
docker push "ioth-registry:5000/registry-update-certs:latest"

docker service rm "ioth-registry-update-certs"

docker service create \
  --detach \
  --name="ioth-registry-update-certs" \
  --restart-condition=none \
  --mount type=bind,source=/etc/hosts,target=/etc/hosts \
  --mount type=bind,source=/etc/docker,target=/etc/docker \
  --mode global \
  "ioth-registry:5000/registry-update-certs:latest" "${BIND_ADDRESS}"
