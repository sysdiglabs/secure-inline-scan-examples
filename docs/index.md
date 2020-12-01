---
title: Sysdig Secure Inline Scan Examples
summary: >
  This is not a comprehensive catalog of examples for all integrations available, but a live document where we continually publish more information as we see users need it.
  We do try to keep a list of links to all integrations and other related websites that you may find useful.
---

# Common scenarios & recipes

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

Assuming the image <image-name> is avaiable as an image tarball at `image.tar`.

For example, the command `docker save <image-name> -o image.tar` creates a tarball for <image-name>.

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

# Other integrations and examples

In this [repository](https://github.com/sysdiglabs/secure-inline-scan-examples/) you can find the following examples in alphabetical order:

* [Google Cloud Build](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/google-cloud-build)
* Jenkins
  * [Scan from repository](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/jenkins-scan-from-repo)
  * [Build and scan](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/jenkins-build-push-scan-from-repo)
  * [Build, push and scan from repository](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/jenkins/jenkins-build-and-scan)
* [Tekton](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/tekton)
  * [Tekton alpha API](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/tekton/alpha)
  * [Tekton beta API](https://github.com/sysdiglabs/secure-inline-scan-examples/tree/main/tekton/beta)
* Unprivileged Docker
  * [Scan from local build](https://github.com/sysdiglabs/secure-inline-scan-examples/blob/main/unprivileged-docker/localbuild_scan.sh)
  * [Scan from registry](https://github.com/sysdiglabs/secure-inline-scan-examples/blob/main/unprivileged-docker/registry_scan.sh)

# Other sources of information

The following content is related to inline scanning, and lives outside this repository.

## Integrations

These integrations have a specific entry in their respective CI/CD catalogs:

  * [Jenkins plugin](https://plugins.jenkins.io/sysdig-secure/)
  * [GitHub Action](https://github.com/marketplace/actions/sysdig-secure-inline-scan)

## Documentation pages

Official documentation pages must be current to the features provided by the inline scanner, but their explanations may be brief:

* [Registry Scanning](https://sysdig.com/products/kubernetes-security/image-scanning/) (main Sysdig web page)
* [Image Scanning](https://docs.sysdig.com/en/image-scanning.html) (Sysdig Documentation website)
* [Sysdig Secure inline scan repository](https://github.com/sysdiglabs/secure-inline-scan) (main project code repository's readme)

## Blog articles

Blog articles contain detailed step by step information, but may be out of date respect their current implementations:

* [Image scanning for Google Cloud Build](https://sysdig.com/blog/image-scanning-google-cloud-build/) <nobr>📅 2020-10-06</nobr>
* [Automate Fargate image scanning](https://sysdig.com/blog/fargate-image-scanning/) <nobr>📅 2020-09-29</nobr>
* [Automate registry scanning with Harbor & Sysdig](https://sysdig.com/blog/harbor-registry-scanning/) <nobr>📅 2020-08-11</nobr>
* [12 Container image scanning best practices to adopt in production](https://sysdig.com/blog/image-scanning-best-practices/) <nobr>📅 2020-07-21</nobr>
* [Performing Image Scanning on Admission Controller with OPA](https://sysdig.com/blog/image-scanning-admission-controller/) <nobr>📅 2020-04-16</nobr>
* [Securing Tekton pipelines in OpenShift with Sysdig](https://sysdig.com/blog/securing-tekton-pipelines-openshift/) <nobr>📅 2020-04-09</nobr>
* [Image scanning for CircleCI](https://sysdig.com/blog/image-scanning-circleci/) <nobr>📅 2020-02-20</nobr>
* [Image scanning for Gitlab CI/CD](https://sysdig.com/blog/gitlab-ci-cd-image-scanning/) <nobr>📅 2020-01-26</nobr>
* [Image Scanning with Github Actions](https://sysdig.com/blog/image-scanning-github-actions/) <nobr>📅 2020-01-14</nobr>
* [AWS ECR Scanning with Sysdig Secure](https://sysdig.com/blog/aws-ecr-scanning/) <nobr>📅 2019-11-26</nobr>
* [Inline Image Scanning for AWS CodePipeline and AWS CodeBuild](https://sysdig.com/blog/image-scanning-aws-codepipeline-codebuild/) <nobr>📅 2019-11-26</nobr>
* [Image scanning for Azure Pipelines](https://sysdig.com/blog/image-scanning-azure-pipelines/) <nobr>📅 2019-10-29
* [Docker scanning for Jenkins CI/CD security with the Sysdig Secure plugin](https://sysdig.com/blog/docker-scanning-jenkins/) <nobr>📅 2018-09-05</nobr>
* [Scanning images in Azure Container Registry](https://sysdig.com/blog/scanning-images-in-azure-container-registry/) <nobr>📅 2018-09-04</nobr>

# Contributing

If you find a related topic lacks enough information, or some problem with any of the existing examples, please file a issue in this [repository](https://github.com/sysdiglabs/secure-inline-scan-examples/). Pull requests to ammend any existing information or examples are also welcomed.
