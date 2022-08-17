# GitLab CI Demo - No DinD

In this demo we will use GitLab pipelines without requiring privileged containers, or docker in docker.
We will need to split this pipeline into three different jobs
1. Kaniko: Tool used to build docker image
2. Sysdig-cli-scanner: Scan docker images for vulnerabilities using the new scan engine developed by Sysding in 2022
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
The scan stage is using `sysdig-cli-scanner`. This stage uses a the latest Sysdig scanning method documented here [Sysdig Secure - Vulnerabilities](https://docs.sysdig.com/en/docs/sysdig-secure/vulnerabilities/pipeline/)
We then save the `build/` directory as an artifact for the next step as well as the `report/` directory to review the PDF scan results later.

### Push
The push stage is using `crane`. It simply authenticates to your docker registry and pushes the conatiner from the Build stage to the remote registry
