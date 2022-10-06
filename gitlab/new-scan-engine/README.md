# GitLab CI Demo

In this demo we will use GitLab CI/CD pipelines. We will need to split this pipeline into three different jobs:

1. Kaniko: Tool used to build docker image
2. Sysdig-cli-scanner: Scan docker images for vulnerabilities using the new scan engine developed by Sysding in 2022
3. Crane: Push container image to a remote registry

The pipeline leverages the GitLab's container registry to store the container image once the scan has been successfully completed. There are a few special CI/CD variables to use the Container registry (`CI_REGISTRY*`) that are populated automatically by GitLab so there is no need to specify them in our pipeline if we want to use it, cool!

The [official documentation](https://docs.gitlab.com/ee/user/packages/container_registry/index.html#authenticate-by-using-gitlab-cicd) explains this in more detail but the following is an example of the variables' content once they are [automatically populated](https://docs.gitlab.com/ee/ci/variables/#list-all-environment-variables):

```
CI_REGISTRY="registry.example.com"
CI_REGISTRY_IMAGE="registry.example.com/gitlab-org/gitlab-foss"
CI_REGISTRY_USER="gitlab-ci-token"
CI_REGISTRY_PASSWORD="[masked]"
```

## Setup

In the GitLab repo settings add the `SYSDIG_SECURE_TOKEN` variable to store the Sysdig Token.

Modify the `gitlab-ci.yml` file to replace the image tag if needed:

```
CI_IMAGE_TAG: "latest"
```

## Pipeline stages

### Build

The build stage leverages Kaniko. The container is built as an OCI format tarball file in `$(pwd)/build/$CI_IMAGE_TAG.tar` and not pushed to a remote registry (it will be done only if the scan is successful).

### Scan

The scan stage leverages `sysdig-cli-scanner`. This stage uses the latest Sysdig scanning method documented here [Sysdig Secure - Vulnerabilities](https://docs.sysdig.com/en/docs/sysdig-secure/vulnerabilities/pipeline/).

### Push

The push stage uses `crane` to authenticate to the GitLab registry and to push the container image already built from the Build stage to the remote registry.
