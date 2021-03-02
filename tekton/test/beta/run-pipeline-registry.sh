#!/bin/bash

# Execute example
export DOCKER_USER=$(cat $KEYS/DOCKER_USER)
cat ../../beta/tekton-inline-scan-registry-beta.yaml | envsubst | \
    kubectl create -n tekton-pipelines -f -
