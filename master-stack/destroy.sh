#!/bin/bash

SERVICE="$1"

if [ $# -eq 0 ]; then
    echo "No arguments supplied"
    exit 1
fi

if [ -f "docker-compose.${SERVICE}.yml" ]; then
    docker-compose -p ioth_${SERVICE} -f docker-compose.${SERVICE}.yml down
fi

if [ -f "docker-stack.${SERVICE}.yml" ]; then
    docker stack rm ioth_${SERVICE}
fi
