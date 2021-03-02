#!/bin/bash
stern tekton-pipelines-controller --exclude 'level.:.info' --exclude 'Failed to log the metrics'