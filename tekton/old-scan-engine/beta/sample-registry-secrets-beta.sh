#!/bin/bash
set -x

kubectl create secret docker-registry regcred \
                    --docker-server=index.docker.io \
                    --docker-username=$DOCKER_USER \
                    --docker-password=$DOCKER_PASS \
                    --docker-email=$DOCKER_EMAIL
                     -n tekton-pipelines