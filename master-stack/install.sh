#!/bin/bash
set -e

read -ep "Enter your main domain (e.g. example.com): " DOMAIN
read -ep "Enter your new password (default: auto-generate): " PASSWORD

cp .env.example .env
sed -i "s/example.com/$DOMAIN/g" .env

mkdir ./data

# traefik
cp -a ./template/traefik ./data/traefik

mkdir ./data/traefik/logs
touch ./data/traefik/acme.json
chmod 700 ./data/traefik/acme.json

# portainer
mkdir ./data/portainer

# authelia
cp -a ./template/authelia ./data/authelia
touch ./data/authelia/notification.txt

## secrets
mkdir ./data/authelia/secrets
openssl rand -base64 32 > ./data/authelia/secrets/jwt
openssl rand -base64 32 > ./data/authelia/secrets/session

## generate new password
if [[ $PASSWORD == "" ]]; then
    PASSWORD=$(openssl rand -hex 12)
    echo "Your new password is: $PASSWORD"
fi

## replace domain and password
sed -i "s/example.com/$DOMAIN/g" ./data/authelia/{configuration.yml,users_database.yml}
PASSWORD=$(docker run --rm authelia/authelia authelia hash-password $PASSWORD | sed 's/Password hash: //g')
sed -i "s/<PASSWORD>/$(echo $PASSWORD | sed -e 's/[\/&]/\\&/g')/g" ./data/authelia/users_database.yml

# fluentd
cp -a ./template/fluent ./data/fluent

# elasticsearch
mkdir ./data/elasticsearch

# alertmanager
cp -a ./template/alertmanager ./data/alertmanager
mkdir ./data/alertmanager/data

# prometheus
cp -a ./template/prometheus ./data/prometheus
mkdir ./data/prometheus/data

# grafana
cp -a ./template/grafana ./data/grafana
mkdir ./data/grafana/data
