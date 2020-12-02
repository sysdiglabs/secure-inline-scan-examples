#!/bin/bash

kubectl delete secret regcred -n tekton-pipelines
kubectl delete secret sysdig-secrets -n tekton-pipelines

kubectl delete clusterrole tutorial-role
kubectl delete clusterrolebinding tutorial-binding
