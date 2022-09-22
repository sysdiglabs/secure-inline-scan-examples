#!/bin/bash

echo "Waiting for Tekton pods to be ready..."
kubectl wait --for=condition=ready pod -l 'app.kubernetes.io/part-of=tekton-pipelines'
