#!/bin/bash

crc setup
crc start -p $KEYS/crc-pull-secret.txt

./init/oc-login.sh
./init/openshift-project.sh
