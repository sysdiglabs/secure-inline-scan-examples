#!/bin/bash

cd alpha
./init-tekton-alpha.sh
./wait-tekton-ready.sh
./prepare-credentials.sh
./run-pipeline-registry.sh
