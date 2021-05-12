#!/bin/sh

BIND_ADDRESS="10.8.0.1"
BIND_PORT="7999"

# replace /etc/hosts content
tmpfile=$(mktemp)
cp "/etc/hosts" "${tmpfile}"
sed '/ioth-registry/d' -i "${tmpfile}"
sh -c 'echo "'"${BIND_ADDRESS}"' ioth-registry" >> '"${tmpfile}"
cat "${tmpfile}" > "/etc/hosts"
rm "${tmpfile}"

# apply certificates
mkdir -p "/etc/docker/certs.d/ioth-registry:5000"
wget -O "/etc/docker/certs.d/ioth-registry:5000/$(date '+%Y-%m-%d_%H-%M-%S').crt" "${BIND_ADDRESS}:${BIND_PORT}/domain.crt"
