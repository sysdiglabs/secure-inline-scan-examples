#!/bin/bash

REPO=quay.io/e_minguez/sysdig-cli-scanner
export VERSION=$(curl -L -s https://download.sysdig.com/scanning/sysdig-cli-scanner/latest_version.txt)

docker build --build-arg VERSION . -t ${REPO}:${VERSION}
docker push ${REPO}:${VERSION}