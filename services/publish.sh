#!/bin/bash

BASE="./"
for SERVICE_PATH in $(find $BASE -maxdepth 1 -type d)
do
    if [ "$SERVICE_PATH" = "$BASE" ]
    then
        continue
    fi

    SERVICE_NAME="${SERVICE_PATH#"$BASE"}"
	docker push "ioth-registry:5000/$SERVICE_NAME:latest"
done
