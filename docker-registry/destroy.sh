#!/bin/sh
cd "$(dirname "$0")"

# remove container if already exists
docker stop "ioth-registry"
docker rm "ioth-registry"

# remove from /etc/hosts
sed '/ioth-registry/d' -i /etc/hosts

# remove from certificated
rm -rf "/etc/docker/certs.d/ioth-registry:5000"

# remove generated certs
rm -rf "certs"
