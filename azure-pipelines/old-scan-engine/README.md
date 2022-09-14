# Azure Pipelines Demo

In this demo we will use Azure Pipelines to build, scan and push a container image.

NOTE: This example uses the [legacy Sysdig scanning engine](https://docs.sysdig.com/en/docs/sysdig-secure/scanning/)

The workflow is as follows:

1. Build the container image and store it locally
2. Run the `sysdiglabs/secure-inline-scan:2` container to perform the scan
3. Push the container image to a remote registry

## Setup

### Variables

It is required to create a `secureApiKey` pipeline variable containing the Sysdig API token in order
to be able to perform the scan. See [the official documentation](https://docs.microsoft.com/en-us/azure/devops/pipelines/process/set-secret-variables)
for instructions on how to do it, but basically:

* Edit the pipeline
* Select "Variables"
* Add a new `secureApiKey` variable with the proper content

### Registry access

It is required to create a Docker registry "Service Connections" to be able to push images to the registry.
See [the official documentation](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml#docker-hub-or-others)
for instructions on how to do it, but basically:

* Select Project settings > Service connections
* Select + New service connection, select the "Docker Registry", and then select Next
* Add the registry url, user & password and a Service connection name (in this example, the Service connection name is `containerRegistry`)

Then, modify the variables on the [azure-pipelines.yml](azure-pipelines.yml) file to fit your needs:

```
containerRegistryConnection: containerRegistry
imageName: "sysdiglabs/dummy-vuln-app"
tags: "latest"
```
