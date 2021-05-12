#!/bin/sh

BIND_ADDRESS="${1-10.8.0.1}"
BIND_PORT="${2-7999}"

docker run --rm \
   -v ${PWD}/../certs/domain.crt:/usr/share/nginx/html/domain.crt \
   -v ${PWD}/install.sh:/usr/share/nginx/html/install.sh.template \
   -v ${PWD}/entrypoint.sh:/docker-entrypoint.d/40-replace-ioth-ip.sh \
   -p ${BIND_ADDRESS}:${BIND_PORT}:80 \
   -e "BIND_ADDRESS=${BIND_ADDRESS}" \
   -e "BIND_PORT=${BIND_PORT}" \
   nginx;
