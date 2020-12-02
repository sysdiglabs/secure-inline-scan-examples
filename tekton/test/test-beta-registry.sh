#!/bin/bash

cd beta
./init-tekton-beta.sh
./wait-tekton-ready.sh
./prepare-credentials.sh
./run-pipeline-registry.sh
