# Unsupported container for the `sysdig-cli-scanner`

A few notes:
* It does a multistep build to get the binary and `chmod`-it from an alpine container, then it uses the debian distroless to save some disk space (the binary itself is 28MB and the image is 31MB)
* The `sysdig-cli-scanner` version number is used also for the container label
* The container image itself is scanned by the `sysdig-cli-scanner`!

## Run it

```
$ docker run -e SECURE_API_TOKEN="X" ghcr.io/sysdiglabs/sysdig-cli-scanner:1.2.10 --apiurl https://eu1.app.sysdig.com pull://docker.io/sysdiglabs/dummy-vuln-app
2022-10-28T10:23:05Z Starting analysis with Sysdig scanner version 1.2.10-rc
2022-10-28T10:23:05Z Retrieving vulnerabilities DB...
2022-10-28T10:23:07Z Done 116.3 MB
2022-10-28T10:23:07Z Loading vulnerabilities DB...
2022-10-28T10:23:07Z Done
2022-10-28T10:23:07Z Retrieving image...
2022-10-28T10:23:08Z Done
2022-10-28T10:23:08Z Scan started...
2022-10-28T10:23:16Z Uploading result to backend...
2022-10-28T10:23:16Z Done
2022-10-28T10:23:16Z Total execution time 11.019413828s

Type: dockerImage
ImageID: sha256:b670c067178c876d17363baec279d483ae07384351d1a0be7646230442471ac6
Digest: sysdiglabs/dummy-vuln-app@sha256:bc86e8ba5741ab71ce50f13fbf89a1f27dc4e1d3b0c3345cee8e3238bc30022b
BaseOS: debian 9.13
PullString: docker.io/sysdiglabs/dummy-vuln-app

13 vulnerabilities found
2 Critical (0 fixable)
5 High (2 fixable)
6 Medium (5 fixable)
0 Low (0 fixable)
0 Negligible (0 fixable)

  PACKAGE   TYPE   VERSION  SUGGESTED FIX  CRITICAL  HIGH  MEDIUM  LOW  NEGLIGIBLE  EXPLOIT
  pip      python   9.0.1       19.2          0       2      1      0       0          0
  numpy    python  1.12.1      1.19.0         0       1      3      0       0          0
  pyxdg    python   0.25        0.26          0       1      0      0       0          0
  Jinja2   python  2.11.2      2.11.3         0       0      1      0       0          0

                  POLICIES EVALUATION
    Policy: Sysdig Best Practices FAILED (8 failures)

Policies evaluation FAILED at 2022-10-28T10:23:16Z
Full image results here: https://eu1.app.sysdig.com/secure/#/scanning/assets/results/1722348e04906294017718c0cd082970/overview (id 1722348e04906294017718c0cd082970)
Execution logs written to: /home/nonroot/scan-logs
```

## Build it

The container is built by the [GitHub workflow](../.github/workflows/sysdig-cli-scanner.yaml) but in order to do it manually you can use the [doit.sh](./doit.sh) script. It requires you to be logged in your container image repository (docker login) and modify the REPO variable in the doit.sh script.
