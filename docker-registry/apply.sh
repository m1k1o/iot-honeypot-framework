#!/bin/sh
cd "$(dirname "$0")"

BIND_ADDRESS="${1}"

# replace /etc/hosts content
tmpfile=$(mktemp)
cp "/etc/hosts" "${tmpfile}"
sed '/ioth-registry/d' -i "${tmpfile}"
sh -c 'echo "'"${BIND_ADDRESS}"' ioth-registry" >> '"${tmpfile}"
cat "${tmpfile}" > "/etc/hosts"
rm "${tmpfile}"

# apply certificates
mkdir -p "/etc/docker/certs.d/ioth-registry:5000"
cp "certs/domain.crt" "/etc/docker/certs.d/ioth-registry:5000/$(date '+%Y-%m-%d_%H-%M-%S').crt"
