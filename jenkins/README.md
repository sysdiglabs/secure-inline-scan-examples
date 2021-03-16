# Inline Scan in Jenkins

## Sysdig Secure Jenkins plugin

The [Sysdig Secure Jenkins plugin](https://plugins.jenkins.io/sysdig-secure/) can be used in a Pipeline job, or added as a build step to a Freestyle job to automate the process of running an image analysis, evaluating custom policies against images, and performing security scans.

The plugin supports both backend and inline scanning and scan result integration. It publishes the scan report as part of the Jenkins job results, directly available from the UI.

However, the current version requires a working Docker environment (Docker socket must be available) for inline scanning. This requirement doesn't fit all scenarios, like running the Jenkins worker as a Pod when using the Kubernetes plugin, unless you use [Docker-in-Docker, which is discouraged](https://jpetazzo.github.io/2015/09/03/do-not-use-docker-in-docker-for-ci/).

See more information at the plugin page: https://plugins.jenkins.io/sysdig-secure/

## Using the Kubernetes plugin (podTemplate)

The inline scanner runs as a container `quay.io/sysdig/secure-inline-scan:2`, but it does not depend on the Docker socket being available (except for scanning local Docker images). So the inline scan container can be executed as part of a [Kubernetes plugin](https://plugins.jenkins.io/kubernetes/) podTemplate.

The following pipelines show different usage examples. All they have in common is the `quay.io/sysdig/secure-inline:2` container is added as an additional container inside the podTemplate executing the Jenkins worker. The container entrypoint is changed to `cat` so the container is started in a "paused" state. At some point in the pipeline, the `/sysdig-inline-scan.sh` script (the original entrypoint)  is executed inside the inline-scan container.

  * [Scan from repository](jenkins-scan-from-repo/)
  * [Build and scan](jenkins-build-and-scan/)
  * [Build, push and scan from repository](jenkins-build-push-scan-from-repo/)
  * [Build, push and scan using Openshift internal registry](jenkins-openshift-internal-registry/)

### Troubleshooting

#### Execution of `/sysdig-inline-scan.sh` getting stuck 

In case the execution of the:

```
                container("inline-scan") {
                    sh "/sysdig-inline-scan.sh -k ${SECURE_API_KEY_PSW} ${IMAGE_NAME}"
                }
```

is stuck for a while, and then finishes with:

```
...
process apparently never started in /tmp/workspace/test-sysdig-inline-scan@tmp/durable-c48f0134
(running Jenkins temporarily with -Dorg.jenkinsci.plugins.durabletask.BourneShellScript.LAUNCH_DIAGNOSTICS=true might make the problem clearer)
```

you can enable launch diagnostics as described in the error message to get further information. You might find one of the following problems.

**UID mismatch**

Described in the [Kubernetes plugin page](https://plugins.jenkins.io/kubernetes/):

> this problem usually happens when the UID of the user in the JNLP container differs from the one in other container(s). All containers you use should have the same UID of the user, also this can be achieved by setting *securityContext*

For example, if the *jnlp* container runs as UID 1001, and the inline-scan container runs with the default UID 1000, the `sh` step will fail due to permissions when writing the output:

```
sh: can't create /home/jenkins/agent/workspace/thejob@tmp/durable-e0b7cd27/jenkins-log.txt: Permission denied
sh: can't create /home/jenkins/agent/workspace/thejob@tmp/durable-e0b7cd27/jenkins-result.txt.tmp: Permission denied
mv: can't rename '/home/jenkins/agent/workspace/thejob@tmp/durable-e0b7cd27/jenkins-result.txt.tmp': No such file or directory
touch: /home/jenkins/agent/workspace/thejob@tmp/durable-e0b7cd27/jenkins-log.txt: Permission denied
touch: /home/jenkins/agent/workspace/thejob@tmp/durable-e0b7cd27/jenkins-log.txt: Permission denied
touch: /home/jenkins/agent/workspace/thejob@tmp/durable-e0b7cd27/jenkins-log.txt: Permission denied
```

You can fix it by making all the containers execute using the same UID using the [securityContext parameter](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) in the podTemplate

**Working directory**

If the `workingDir` of the container is not specified, it could be set to an invalid path that causes problems when mounting. Try setting `workingDir: /tmp`.
