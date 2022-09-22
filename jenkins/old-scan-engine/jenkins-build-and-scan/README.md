# Build and scan example

This [example pipeline](Jenkinsfile) shows how to build, scan and push a Docker image in a dockerless environment, by creating a podTemplate with 4 containers:
 * **jnlp** container. Required for the Jenkins agent.
 * **maven** container for building a Java application.
 * **builder** container, using [Kaniko](https://github.com/GoogleContainerTools/kaniko) to build a Docker image without requiring the Docker daemon. The `--no-push` option tells Kaniko not to push the image to a registry, and just store it locally in the `oci` folder inside the workspace.
 * **inline-scan** container, where the pipeline executes the `inline-scan.sh` script to analyze the image built in the previous step locally from the `oci` folder in the workspace.

Finally, an optional 5th step could be included to push the image to the registry, if the scan is successful, by using [Skopeo](https://github.com/containers/skopeo).

See [Jenkins examples README.md](../README.md) for common usage tips and troubleshooting.