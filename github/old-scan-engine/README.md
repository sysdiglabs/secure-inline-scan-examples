# GitHub CI Demo

In this demo we will use GitHub actions to build, scan and push a container image.

The workflow is based on the [sysdiglabs/dummy-vuln-app](https://github.com/sysdiglabs/dummy-vuln-app) application and and uses the [Sysdiglabs/scan-action](https://github.com/sysdiglabs/scan-action) GitHub action to scan it.

The workflow is as follows:

1. Build the container image and store it locally
2. Perform the scan using the [Sysdiglabs/scan-action](https://github.com/sysdiglabs/scan-action)
3. Upload a SARIF report

## Setup

It is required to create a repository secret to store the Sysdig Token:

* `SYSDIG_SECURE_TOKEN`: Sysdig Token
