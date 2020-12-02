#!/bin/bash

export DOCKER_USER=$(cat $KEYS/DOCKER_USER)
cat ../../beta/tekton-inline-scan-localbuild-beta.yaml | envsubst | \
    kubectl create -n tekton-pipelines -f -
