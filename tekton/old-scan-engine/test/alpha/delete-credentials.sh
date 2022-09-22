#!/bin/bash

kubectl delete secret docker-auth-for-tekton -n tekton-pipelines
kubectl delete secret sysdig-secrets -n tekton-pipelines
