#!/bin/bash
CRC_KUBEADMIN_KEY=$(cat $KEYS/CRC_KUBEADMIN_KEY)

export PATH=$HOME/.crc/bin:$PATH
oc login -u kubeadmin -p $CRC_KUBEADMIN_KEY https://api.crc.testing:6443

oc adm policy add-scc-to-user anyuid -z tekton-pipelines-controller