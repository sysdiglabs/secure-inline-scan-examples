#!/bin/bash

export PATH=$HOME/.crc/bin:$PATH
oc new-project tekton-pipelines
oc adm policy add-scc-to-user anyuid -z tekton-pipelines-controller