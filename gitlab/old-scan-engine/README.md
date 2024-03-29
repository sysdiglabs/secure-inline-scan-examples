# GitLab CI Demo - No DinD

> :warning: **Outdated example**: This example is using the legacy scan engine. Please use the [latest example for the new scan engine](../new-scan-engine/README.md) instead.

![Gitlab job](gitlab.png)

In this demo we will use GitLab pipelines without requiring privileged containers, or docker in docker.
We will need to split this pipeline into three different jobs
1. Kaniko: Tool used to build docker image
2. Sysdig-inline-scan (deprecated): Scan docker images for vulnerabilities
3. Crane: Push container image to a remote registry

## Setup
In GitLab repo settings add variables
`CI_REGISTRY_USER`: Docker username
`CI_REGISTRY_PASSWORD`: Docker user password
`SYSDIG_SECURE_TOKEN`: Sysdig Token

Modify the gitlab-ci.yml file to build the image
```
  CI_REGISTRY_HOST: "docker.io"
  CI_REGISTRY_NAME: my-registry
  CI_IMAGE_NAME: "my-image"
  CI_IMAGE_TAG: "latest"
```

The variables are to build the full image url
`$CI_REGISTRY_HOST/$CI_REGISTRY_NAME/$CI_IMAGE_NAME:$CI_IMAGE_TAG`
We would expect
`docker.io/my-registry/my-image:latest`

## Understanding the stages
In order to get around using Docker in docker, these additional stages are necessary

There are three pipeline stages
1. Build
2. Scan
3. Push

### Build
The build stage is using Kaniko. We use a method to build the container to an oci format tarball, saved to the current working directory in `build/` directory. It is not pushed to a remote registry.
We then save the `build/` directory as an artifact.

### Scan
The scan stage is using `sysdig-inline-scan:2` (deprecated). This stage uses a scanning method without the docker daemon dependencies ([Documentation](https://docs.sysdig.com/en/docs/sysdig-secure/scanning/integrate-with-cicd-tools/)).
We then save the `build/` directory as an artifact for the next step as well as the `report/` directory to review the PDF scan results later.

### Push
The push stage is using `crane`. It simply authenticates to your docker registry and pushes the conatiner from the Build stage to the remote registry
