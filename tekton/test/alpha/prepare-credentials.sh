#!/bin/bash

export DOCKER_USER=$(cat $KEYS/DOCKER_USER)
export DOCKER_PASS=$(cat $KEYS/DOCKER_PASS)
export SYSDIG_SECURE_API_TOKEN_BASE64=$(printf "$(cat $KEYS/SYSDIG_SECURE_API_TOKEN)" | base64)

cat ../../alpha/sample-registry-secrets.yaml | envsubst | \
    kubectl apply -n tekton-pipelines -f -

cat ../../alpha/sample-sysdig-secrets.yaml | envsubst | \
    kubectl apply -n tekton-pipelines -f -
