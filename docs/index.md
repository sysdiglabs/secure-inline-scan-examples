---
title: Sysdig Vulnerability Scan Examples
summary: >
  This is not a comprehensive catalog of examples for all integrations available, but a live document where we continually publish more information as we see users need it.
  We do try to keep a list of links to all integrations and other related websites that you may find useful.
---

# Legacy Scanner engine vs Vulnerability Management engine

**As of April 20, 2022, Sysdig offers both a Legacy Scanner engine and the newer Vulnerability Management engine. See the [official documentation](https://docs.sysdig.com/en/docs/sysdig-secure/scanning/new-scanning-engine/#which-engine-is-enabled-now) to understand which engine is enabled into your account.**

- [Vulnerability Management engine common scenarios & recipes](#vulnerability-management-engine-common-scenarios--recipes)
  * [Download the `sysdig-cli-scanner`](#download-the-sysdig-cli-scanner)
  * [Scan local image, built using docker](#scan-local-image-built-using-docker)
  * [Local image (provided docker archive)](#local-image-provided-docker-archive)
  * [Public registry image](#public-registry-image)
  * [Private registry image](#private-registry-image)
  * [Containers-storage (cri-o, podman, buildah and others)](#containers-storage-cri-o-podman-buildah-and-others)
- [Legacy Scanner engine common scenarios & recipes](#legacy-scanner-engine-common-scenarios--recipes)
  * [Scan local image, built using docker](#scan-local-image-built-using-docker-1)
  * [Local image (provided docker archive)](#local-image-provided-docker-archive-1)
  * [Public registry image](#public-registry-image-1)
  * [Private registry image](#private-registry-image-1)
  * [Containers-storage (cri-o, podman, buildah and others)](#containers-storage-cri-o-podman-buildah-and-others-1)
  * [Using a proxy](#using-a-proxy)
- [Other integrations and examples](#other-integrations-and-examples)
  * [Vulneratbility Management Engine (new scan engine)](#vulneratbility-management-engine-new-scan-engine)
  * [Legacy Scanner Engine (old scan engine)](#legacy-scanner-engine-old-scan-engine)
- [Other sources of information](#other-sources-of-information)
  * [Integrations](#integrations)
  * [Documentation pages](#documentation-pages)
  * [Blog articles](#blog-articles)
- [Contributing](#contributing)

# Vulnerability Management engine common scenarios & recipes

## Download the `sysdig-cli-scanner`

Linux or MacOS:

```
curl -LO "https://download.sysdig.com/scanning/bin/sysdig-cli-scanner/$(curl -L -s https://download.sysdig.com/scanning/sysdig-cli-scanner/latest_version.txt)/$(uname -s | tr '[:upper:]' '[:lower:]')/amd64/sysdig-cli-scanner"
```

Set the executable flag on the file:

```
chmod +x ./sysdig-cli-scanner
```

You only need to download and set executable once. Then you can scan images by running the `sysdig-cli-scanner` command:

```
SECURE_API_TOKEN=<your-api-token> ./sysdig-cli-scanner --apiurl <sysdig-api-url> <image-name>
```

## Scan local image, built using docker

```
# Build the image locally
docker build -t <image-name> .

# Scan the image, available on local docker
SECURE_API_TOKEN=<your-api-token> ./sysdig-cli-scanner --apiurl <sysdig-api-url> docker://<image-name>
```

## Local image (provided docker archive)

Assuming the image `<image-name>` is available as an image tarball at `image.tar`.

For example, the command `docker save <image-name> -o image.tar` creates a tarball for `<image-name>`.

```
SECURE_API_TOKEN=<your-api-token> ./sysdig-cli-scanner --apiurl <sysdig-api-url> file://tmp/image.tar
```

## Public registry image

Example: scan `alpine` image from public registry. The scanner will pull and scan it.

```
SECURE_API_TOKEN=<your-api-token> ./sysdig-cli-scanner --apiurl <sysdig-api-url> pull://alpine
```

## Private registry image

To scan images from private registries, you might need to provide credentials:

```
$ REGISTRY_USER=<YOUR_REGISTRY_USERNAME> REGISTRY_PASSWORD=<YOUR_REGISTRY_PASSWORD> SECURE_API_TOKEN=<YOUR_API_TOKEN> ./sysdig-cli-scanner --apiurl https://secure.sysdig.com ${REPO_NAME}/${IMAGE_NAME}
```

## Containers-storage (cri-o, podman, buildah and others)

Scan images from container runtimes using containers-storage format:

```
# Build an image using buildah from a Dockerfile
buildah build-using-dockerfile -t myimage:latest

# Scan the image
SECURE_API_TOKEN=<your-api-token> ./sysdig-cli-scanner --apiurl <sysdig-api-url> crio://localhost/myimage:latest
```

Example for an image pulled with podman

```
podman pull docker.io/library/alpine

#Scan the image
SECURE_API_TOKEN=<your-api-token> ./sysdig-cli-scanner --apiurl <sysdig-api-url> podman://docker.io/library/alpine
```

# Legacy Scanner engine common scenarios & recipes

## Scan local image, built using docker

```
#Build the image locally
docker build -t <image-name> .

#Scan the image, available on local docker. Mounting docker socket is required
docker run --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    quay.io/sysdig/secure-inline-scan:2 \
    --sysdig-url <omitted> \
    --sysdig-token <omitted> \
    --storage-type docker-daemon \
    --storage-path /var/run/docker.sock \
    <image-name>
```

## Local image (provided docker archive)

Assuming the image `<image-name>` is available as an image tarball at `image.tar`.

For example, the command `docker save <image-name> -o image.tar` creates a tarball for `<image-name>`.

```
docker run --rm \
    -v ${PWD}/image.tar:/tmp/image.tar \
    quay.io/sysdig/secure-inline-scan:2 \
    --sysdig-url <omitted> \
    --sysdig-token <omitted> \
    --storage-type docker-archive \
    --storage-path /tmp/image.tar \
    <image-name>
```

## Public registry image

Example: scan `alpine` image from public registry. The scanner will pull and scan it.

```
docker run --rm \
    quay.io/sysdig/secure-inline-scan:2 \
    --sysdig-url <omitted> \
    --sysdig-token <omitted> \
    alpine
```

## Private registry image

To scan images from private registries, you might need to provide credentials:

```
docker run --rm \
    quay.io/sysdig/secure-inline-scan:2 \
    --sysdig-url <omitted> \
    --sysdig-token <omitted> \
    --registry-auth-basic <user:passw> \
    <image-name>
```

Authentication methods available are:
* `--registry-auth-basic` for authenticating via http basic auth
* `--registry-auth-file` for authenticating via docker/skopeo credentials file
* `--registry-auth-token` for authenticating via registry token

## Containers-storage (cri-o, podman, buildah and others)

Scan images from container runtimes using containers-storage format:

```
#Build an image using buildah from a Dockerfile
buildah build-using-dockerfile -t myimage:latest

#Scan the image. Options '-u root' and '--privileged' might be needed depending
#on the access permissions for /var/lib/containers
docker run \
    -u root --privileged \
    -v /var/lib/containers:/var/lib/containers \
    quay.io/sysdig/secure-inline-scan:2 \
    --storage-type cri-o \
    --sysdig-token <omitted> \
    localhost/myimage:latest
```

Example for an image pulled with podman

```
podman pull docker.io/library/alpine

#Scan the image. Options '-u root' and '--privileged' might be needed depending
#on the access permissions for /var/lib/containers
docker run \
    -u root --privileged \
    -v /var/lib/containers:/var/lib/containers \
    quay.io/sysdig/secure-inline-scan:2 \
    --storage-type cri-o \
    --sysdig-token <omitted> \
    docker.io/library/alpine
```

## Using a proxy

To use a proxy, set the standard `http_proxy` and `https_proxy` variables when running the container.

Example:

```
docker run --rm \
    -e http_proxy="http://my-proxy:3128" \
    -e https_proxy="http://my-proxy:3128" \
    quay.io/sysdig/secure-inline-scan:2 \
    --sysdig-url <omitted> \
    --sysdig-token <omitted> \
    alpine
```

Both `http_proxy` and `https_proxy` variables are required, as some tools will use per-scheme proxy.

The `no_proxy` variable can be used to define a list of hosts that don't use the proxy.

# Other integrations and examples

In this [repository](https://github.com/sysdiglabs/secure-inline-scan-examples/) you can find the following examples in alphabetical order:

## Vulneratbility Management Engine (new scan engine)

* [AWS Codebuild](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/aws-codebuild/new-scan-engine)
* [Azure Pipelines](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/azure-pipelines/new-scan-engine)
* [GitLab](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/gitlab/new-scan-engine)
* [GitHub](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/github/new-scan-engine)
* [Jenkins](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/new-scan-engine)


## Legacy Scanner Engine (old scan engine)

* [Azure Pipelines](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/azure-pipelines/old-scan-engine)
* [GitLab](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/gitlab/old-scan-engine)
* [GitHub](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/github/old-scan-engine)
* [Google Cloud Build](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/google-cloud-build/old-scan-engine)
* [Jenkins](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/old-scan-engine)
  * [Scan from repository](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/old-scan-engine/jenkins-scan-from-repo)
  * [Build and scan](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/old-scan-engine/jenkins-build-and-scan)
  * [Build, push and scan from repository](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/old-scan-engine/jenkins-build-push-scan-from-repo)
  * [Build, push and scan using Openshift internal registry](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/old-scan-engine/jenkins-openshift-internal-registry)
* [Tekton](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/tekton/old-scan-engine)
  * [Tekton alpha API](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/tekton/old-scan-engine/alpha)
  * [Tekton beta API](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/tekton/old-scan-engine/beta)
* [Unprivileged Docker](https://github.com/sysdiglabs/secure-inline-scan-examples/blob/main/unprivileged-docker/old-scan-engine)
  * [Scan from local build](https://github.com/sysdiglabs/secure-inline-scan-examples/blob/main/unprivileged-docker/old-scan-engine/localbuild_scan.sh)
  * [Scan from registry](https://github.com/sysdiglabs/secure-inline-scan-examples/blob/main/unprivileged-docker/old-scan-engine/registry_scan.sh)

# Other sources of information

## Integrations

These integrations have a specific entry in their respective CI/CD catalogs:

  * [Jenkins plugin (both new and old scan engines)](https://plugins.jenkins.io/sysdig-secure/)
  * [GitHub Action (old scan engine)](https://github.com/marketplace/actions/sysdig-secure-inline-scan)

## Documentation pages

* [Sysdig - Vulnerability Management](https://docs.sysdig.com/en/docs/sysdig-secure/vulnerabilities/)
* [Sysdig - Scanning (Legacy)](https://docs.sysdig.com/en/docs/sysdig-secure/scanning/)

## Blog articles

Blog articles contain detailed step by step information, but may be out of date respect their current implementations:

* [Image scanning for Google Cloud Build](https://sysdig.com/blog/image-scanning-google-cloud-build/) <nobr>📅 2020-10-06</nobr>
* [Automate Fargate image scanning](https://sysdig.com/blog/fargate-image-scanning/) <nobr>📅 2020-09-29</nobr>
* [Automate registry scanning with Harbor & Sysdig](https://sysdig.com/blog/harbor-registry-scanning/) <nobr>📅 2020-08-11</nobr>
* [12 Container image scanning best practices to adopt in production](https://sysdig.com/blog/image-scanning-best-practices/) <nobr>📅 2020-07-21</nobr>
* [Shielding your Kubernetes runtime with image scanning on admission controller](https://sysdig.com/blog/image-scanning-admission-controller/) <nobr>📅 2021-02-18</nobr>
* [Securing Tekton pipelines in OpenShift with Sysdig](https://sysdig.com/blog/securing-tekton-pipelines-openshift/) <nobr>📅 2020-04-09</nobr>
* [Image scanning for CircleCI](https://sysdig.com/blog/image-scanning-circleci/) <nobr>📅 2020-02-20</nobr>
* [Image scanning for Gitlab CI/CD](https://sysdig.com/blog/gitlab-ci-cd-image-scanning/) <nobr>📅 2022-10-12</nobr>
* [Image Scanning with Github Actions](https://sysdig.com/blog/image-scanning-github-actions/) <nobr>📅 2022-09-26</nobr>
* [AWS ECR Scanning with Sysdig Secure](https://sysdig.com/blog/aws-ecr-scanning/) <nobr>📅 2021-11-23</nobr>
* [Inline Image Scanning for AWS CodePipeline and AWS CodeBuild](https://sysdig.com/blog/image-scanning-aws-codepipeline-codebuild/) <nobr>📅 2019-11-26</nobr>
* [Image scanning for Azure Pipelines](https://sysdig.com/blog/image-scanning-azure-pipelines/) <nobr>📅 2022-09-19</nobr>
* [Container Image Scanning on Jenkins with Sysdig](https://sysdig.com/blog/docker-scanning-jenkins/) <nobr>📅 2022-10-26</nobr>

# Contributing

If you find a related topic lacks enough information, or some problem with any of the existing examples, please file a issue in this [repository](https://github.com/sysdiglabs/secure-inline-scan-examples/). Pull requests to ammend any existing information or examples are also welcomed.
