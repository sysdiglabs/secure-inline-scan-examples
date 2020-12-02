#!/bin/bash

tekton_version=v0.16.3
dashboard_version=v0.9.0

kubectl delete -f https://github.com/tektoncd/dashboard/releases/download/$dashboard_version/tekton-dashboard-release.yaml
kubectl delete -f https://github.com/tektoncd/pipeline/releases/download/$tekton_version/release.notags.yaml
