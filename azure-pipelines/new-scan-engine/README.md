# Azure Pipelines Demo

In this demo we will use Azure Pipelines to build, scan and push a container image.

NOTE: This example uses the [new Sysdig scanning engine](https://docs.sysdig.com/en/docs/sysdig-secure/scanning/new-scanning-engine/)

The workflow is as follows:

1. Build the container image and store it locally
2. Download the sysdig-cli-scanner cli if needed
3. Perform the scan
4. Push the container image to a remote registry

The workflow leverages Azure Pipeline actions cache to avoid downloading the binary,
the databases and the container images if they are available.

## Setup

### Variables

It is required to create a TOKEN pipeline variable containing the Sysdig API token in order
to be able to perform the scan. See [the official documentation](https://docs.microsoft.com/en-us/azure/devops/pipelines/process/set-secret-variables)
for instructions on how to do it, but basically:

* Edit the pipeline
* Select "Variables"
* Add a new TOKEN variable with the proper content

### Registry access

It is required to create a Docker registry "Service Connections" to be able to push images to the registry.
See [the official documentation](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml#docker-hub-or-others)
for instructions on how to do it, but basically:

* Select Project settings > Service connections
* Select + New service connection, select the "Docker Registry", and then select Next
* Add the registry url, user & password and a Service connection name (it will be used as REGISTRY_CONNECTION)

Then, modify the variables on the [azure-pipelines.yml](azure-pipelines.yml) file to fit your needs:

```
SYSDIG_SECURE_ENDPOINT: "https://eu1.app.sysdig.com"
REGISTRY_HOST: "quay.io"
IMAGE_NAME: "e_minguez/my-example-app"
IMAGE_TAG: "latest"
REGISTRY_CONNECTION: "quayio-e_minguez"
```
