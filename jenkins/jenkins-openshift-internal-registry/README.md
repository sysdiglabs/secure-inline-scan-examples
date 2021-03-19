# Build, push and scan from Openshift internal registry

This [example pipeline](Jenkinsfile) shows how to build, push, and then scan the Docker image in Openshift, using the service account credentials to push and scan from the Openshift internal registry.

The podTemplate in the example is composed by 4 containers:
 * **jnlp** container. Required for the Jenkins agent. Also, we mount the service account secret in `/home/jenkins/agent/.dockercfg` to convert the old dockercfg format to the new config.json format required by Kaniko:

```
    sh "echo -n \"{ \\\"auths\\\": \"  > /home/jenkins/agent/config.json"
    sh "cat /home/jenkins/agent/.dockercfg >> /home/jenkins/agent/config.json"
    sh "echo \"}\" >>/home/jenkins/agent/config.json"
```

 * **maven** container for building a Java application.
 * **builder** container, using [Kaniko](https://github.com/GoogleContainerTools/kaniko) to build a Docker image without requiring the Docker daemon. Once build, the image is pushed to the internal Openshift registry, using the credentials at `/home/jenkins/agent/config.json`.
 * **inline-scan** container, where the pipeline executes the `inline-scan.sh` script to analyze the image pushed to the internal Openshift registry, using the credentials from /home/jenkins/agent/config.json or using the .dockercfg file (two alternatives are provided).

See [Jenkins examples README.md](../README.md) for common usage tips and troubleshooting.