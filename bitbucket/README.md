# Bitbucket pipelines demo

In this demo we will use Bitbucket pipelines to build, scan and push a container image.
The workflow is as follows:

1. Download the sysdig-cli-scanner
2. Build the container image and store it locally
3. Perform the scan
4. Login to the registry and push the image

## Setup

It is required to create a few repository secrets in order to be able to push the
container image:

* `DOCKER_USER`: Docker username
* `DOCKER_PASSWORD`: Docker user password
* `DOCKER_REGISTRY`: Docker registry URL
* `SECURE_API_TOKEN`: Sysdig Token
