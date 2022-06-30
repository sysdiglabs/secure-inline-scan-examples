# GitHub CI Demo

In this demo we will use GitHub actions to build, scan and push a container image.
The workflow is as follows:

1. Setup Docker Buildx to be able to build the image
2. Build the container image and store it locally
3. Download the sysdig-cli-scanner cli if needed
4. Perform the scan
5. Login to the registry
6. Push the container image to a remote registry

The workflow leverages GitHub actions cache to avoid downloading the binary or
the databases if they are available.

## Setup

It is required to create a few repository secrets in order to be able to push the
container image:

* `REGISTRY_USER`: Docker username
* `REGISTRY_PASSWORD`: Docker user password
* `SECURE_API_TOKEN`: Sysdig Token

Modify the environment variables on the [build-scan-and-push.yaml](build-scan-and-push.yaml) file to fit your needs:

```
SYSDIG_SECURE_ENDPOINT: "https://secure.sysdig.com"
REGISTRY_HOST: "quay.io"
IMAGE_NAME: "mytestimage"
IMAGE_TAG: "my-tag"
DOCKERFILE_CONTEXT: "github/"
```
