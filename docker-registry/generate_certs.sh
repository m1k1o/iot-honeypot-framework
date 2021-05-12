#!/bin/sh
cd "$(dirname "$0")"

rm -rf "certs"
mkdir -p "certs"

openssl req \
  -newkey rsa:4096 -nodes -sha256 -keyout "certs/domain.key" \
  -x509 -days 365 -out "certs/domain.crt" \
  -subj "/C=SK/ST=Slovakia/L=Bratislava/O=STU/OU=FIIT/CN=ioth-registry"
