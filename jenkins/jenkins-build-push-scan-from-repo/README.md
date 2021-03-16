# Build, push and scan from repository example

This [example pipeline](Jenkinsfile) shows how to build, push, and then scan the Docker image in a dockerless environment, by creating a podTemplate with 4 containers:
 * **jnlp** container. Required for the Jenkins agent.
 * **maven** container for building a Java application.
 * **builder** container, using [Kaniko](https://github.com/GoogleContainerTools/kaniko) to build a Docker image without requiring the Docker daemon. Once build, the image is pushed to the destination target registry and repository.
 * **inline-scan** container, where the pipeline executes the `inline-scan.sh` script to analyze the image pushed to the repository in the previous step.

See [Jenkins examples README.md](../README.md) for common usage tips and troubleshooting.