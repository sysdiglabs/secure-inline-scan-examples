#!/bin/bash

export DOCKER_USER=$(cat $KEYS/DOCKER_USER)

cat ../../beta/tekton-inline-scan-registry-beta.yaml | envsubst | \
    kubectl delete -n tekton-pipelines -f -
