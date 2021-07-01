# Run inline-scan and convert to other formats

This [example script](run-inline-scan.sh) shows how to execute the inline-scan container using the `--format=JSON` flag, and then convert the vulnerability report to CSV and HTML using mustache tempaltes.

It mounts the docker socket at /var/run/docker.sock and scans an image locally available in the Docker daemon. So you need to either build the image, or pull it (`docker pull <image-to-scan>).

You might need to run as root (i.e. using `sudo`) or adjust the docker.sock permissions.

How to use:

* Set SECURE_API_TOKEN environment variable with your Sysdig token value
* Execute `./run-inline-scan-sh <image-to-scan>`