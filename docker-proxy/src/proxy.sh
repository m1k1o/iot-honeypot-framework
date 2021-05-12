#!/bin/bash

exec 1>/proc/1/fd/1
exec 2>/proc/1/fd/2

PROTOCOL="$1" # tcp
HOSTNAME="$2" # example.org
TARGET_PORT="$3" # 80
PUBLISHED_PORT="$4" # 8080

if [ "${USE_NCAT}" == "true" ];
then
	NCAT="/usr/bin/ncat"

	if [ "${PROTOCOL}" == "tcp" ];
	then
		# Simple TCP proxy
		${NCAT} -vlk -p ${PUBLISHED_PORT} -e "${NCAT} -v ${HOSTNAME} ${TARGET_PORT}" 2>&1 | grep --line-buffered Connect
	else
		# Simple UDP proxy
		${NCAT} -vlku -p ${PUBLISHED_PORT} -e "${NCAT} -vu ${HOSTNAME} ${TARGET_PORT}" 2>&1 | grep --line-buffered Connect
	fi
else
	if [ "${PROTOCOL}" == "tcp" ];
	then
		# Simple TCP proxy
		docker-proxy -l :${PUBLISHED_PORT} -r ${HOSTNAME}:${TARGET_PORT}
	else
		# Simple UDP proxy
		docker-proxy -u -l :${PUBLISHED_PORT} -r ${HOSTNAME}:${TARGET_PORT}
	fi
fi
