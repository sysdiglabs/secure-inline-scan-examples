#!/bin/bash

export DOCKER_USER=$(cat $KEYS/DOCKER_USER)

cat ../../alpha/tekton-inline-scan-localbuild-alpha.yaml | envsubst | \
    kubectl apply -n tekton-pipelines -f -
