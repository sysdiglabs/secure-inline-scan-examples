#!/bin/bash

./init/wait-tekton-ready.sh
kubectl port-forward svc/tekton-dashboard -n tekton-pipelines 9097:9097
