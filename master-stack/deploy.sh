#!/bin/bash

SERVICE="$1"

if [ $# -eq 0 ]; then
    echo "No arguments supplied"
    exit 1
fi

if [ -f "docker-stack.${SERVICE}.yml" ]; then
    docker-compose -f docker-stack.${SERVICE}.yml config | docker stack deploy -c - ioth_${SERVICE}
fi

if [ -f "docker-compose.${SERVICE}.yml" ]; then
    docker-compose -p ioth_${SERVICE} -f docker-compose.${SERVICE}.yml up -d
fi
