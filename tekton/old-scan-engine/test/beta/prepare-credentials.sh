#!/bin/bash

export DOCKER_USER=$(cat $KEYS/DOCKER_USER)
export DOCKER_PASS=$(cat $KEYS/DOCKER_PASS)
export DOCKER_EMAIL=$(cat $KEYS/DOCKER_EMAIL)
export SYSDIG_SECURE_API_TOKEN_BASE64=$(printf "$(cat $KEYS/SYSDIG_SECURE_API_TOKEN)" | base64)

cat ../../beta/sample-sysdig-secrets.yaml | envsubst | \
    kubectl apply -n tekton-pipelines -f -

../../beta/sample-registry-secrets-beta.sh

