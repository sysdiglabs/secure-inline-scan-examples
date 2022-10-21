# Jenkins image scanning

There are two different approaches if using Jenkins to scan container images for vulnerabilities with Sysdig Secure:

* Using the `sysdig-cli-scanner` binary
* Using the Sysdig Secure Jenkins Plugin

## Using sysdig-cli-scanner

This [example pipeline](Jenkinsfile-sysdig-cli-scanner) shows how to download and execute the new inline scanner to scan an image.

It requires to configure a Jenkins credential `sysdig-secure-api-token` to store the Sysdig Token (as password)

![Screenshot of Jenkins UI](https://github.com/jenkinsci/sysdig-secure-plugin/raw/main/docs/images/SysdigTokenConfiguration.png)

Then the scan is performed by downloading the `sysdig-cli-scanner` tool against the example image.

For a more elaborated example, see the [GitHub](../../github/new-scan-engine/README.md) example.

## Sysdig Secure Jenkins plugin

The [Sysdig Secure Jenkins plugin](https://plugins.jenkins.io/sysdig-secure/) can be used in a Pipeline job, or added as a build step to a Freestyle job to automate the process of running an image analysis, evaluating custom policies against images, and performing security scans.

See more information at the plugin page: https://plugins.jenkins.io/sysdig-secure/

The [example pipeline](Jenkinsfile-jenkins-plugin) shows how to use it to build and scan a container image.

## Prerequisites

Both approaches require a couple of things:

* A valid Sysdig Secure API token.
* Have access to the image storage, either to the local storage where the image was created or to the registry where it is stored.
* The appropriate Sysdig vulnerability scanning endpoint depending on your region, see [the official documentation](https://docs.sysdig.com/en/docs/administration/saas-regions-and-ip-ranges).
