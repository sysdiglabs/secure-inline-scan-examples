# Build, push and scan from Openshift internal registry

This minimalistic [example pipeline](Jenkinsfile) shows how to execute the inline-scan container as part of a podTemplate.

The podTemplate in the example is composed by 2 containers:
 * **jnlp** container. Required for the Jenkins agent.
 * **inline-scan** container, where the pipeline executes the `inline-scan.sh` script to analyze the image from the registry.

See [Jenkins examples README.md](../README.md) for common usage tips and troubleshooting.