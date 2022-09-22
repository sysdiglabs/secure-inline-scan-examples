#!/bin/bash

tekton_version=v0.10.2
dashboard_version=v0.5.1

kubectl create -f https://github.com/tektoncd/pipeline/releases/download/$tekton_version/release.notags.yaml
kubectl create -f https://github.com/tektoncd/dashboard/releases/download/$dashboard_version/tekton-dashboard-release.yaml

