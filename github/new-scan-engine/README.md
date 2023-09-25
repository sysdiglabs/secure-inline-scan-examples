# GitHub CI Demo

In this demo we will use GitHub actions to build, scan and push a container image.
The workflow is as follows:

1. Setup Docker Buildx to be able to build the image
2. Build the container image and store it locally
3. Download the sysdig-cli-scanner cli if needed
4. Perform the scan
5. Login to the registry
6. Push the container image to a remote registry

The workflow leverages GitHub actions cache to avoid downloading the binary or
the databases if they are available.

## Setup

It is required to create a few repository secrets in order to be able to push the
container image:

* `REGISTRY_USER`: Docker username
* `REGISTRY_PASSWORD`: Docker user password
* `SECURE_API_TOKEN`: Sysdig Token

Modify the environment variables on the [build-scan-and-push.yaml](build-scan-and-push.yaml) file to fit your needs:

```
SYSDIG_SECURE_ENDPOINT: "https://secure.sysdig.com"
REGISTRY_HOST: "quay.io"
IMAGE_NAME: "mytestimage"
IMAGE_TAG: "my-tag"
DOCKERFILE_CONTEXT: "github/new-scan-engine/"
```

# Convert to SARIF output

You can use the script sysdig-to-sarif to convert the JSON output of the CLI scanner to SARIF and upload it to Github Security:

```yaml
jobs:
  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      ...

      - name: Setup cache
        uses: actions/cache@v3
        with:
          path: cache
          key: ${{ runner.os }}-cache-${{ hashFiles('**/sysdig-cli-scanner', '**/latest_version.txt', '**/db/main.db.meta.json', '**/scanner-cache/inlineScannerCache.db') }}
          restore-keys: ${{ runner.os }}-cache-

      - name: Download sysdig-cli-scanner if needed
        run:  |
          curl -sLO https://download.sysdig.com/scanning/sysdig-cli-scanner/latest_version.txt
          mkdir -p ${GITHUB_WORKSPACE}/cache/db/
          if [ ! -f ${GITHUB_WORKSPACE}/cache/latest_version.txt ] || [ $(cat ./latest_version.txt) != $(cat ${GITHUB_WORKSPACE}/cache/latest_version.txt) ]; then
            cp ./latest_version.txt ${GITHUB_WORKSPACE}/cache/latest_version.txt
            curl -sL -o ${GITHUB_WORKSPACE}/cache/sysdig-cli-scanner "https://download.sysdig.com/scanning/bin/sysdig-cli-scanner/$(cat ${GITHUB_WORKSPACE}/cache/latest_version.txt)/linux/amd64/sysdig-cli-scanner"
            chmod +x ${GITHUB_WORKSPACE}/cache/sysdig-cli-scanner
          else
            echo "sysdig-cli-scanner latest version already downloaded"
          fi

      - name: Scan the image using sysdig-cli-scanner
        env:
          SECURE_API_TOKEN: ${{ secrets.SECURE_API_TOKEN }}
        run: |
          ${GITHUB_WORKSPACE}/cache/sysdig-cli-scanner \
            --apiurl ${SYSDIG_SECURE_ENDPOINT} \
            --console-log \
            --json-scan-result report.json \
            <put your image name here>

    - name: Scan the image using sysdig-cli-scanner
      env:
        SECURE_API_TOKEN: ${{ secrets.SECURE_API_TOKEN }}
      run: |
        ${GITHUB_WORKSPACE}/cache/sysdig-cli-scanner \
          --apiurl ${SYSDIG_SECURE_ENDPOINT} \
          docker://${REGISTRY_HOST}/${{github.repository_owner}}/${IMAGE_NAME}:${IMAGE_TAG} \
          --console-log \
          --dbpath=${GITHUB_WORKSPACE}/cache/db/ \
          --cachepath=${GITHUB_WORKSPACE}/cache/scanner-cache/

      - name: Generate SARIF report
        run: |
          python3 /path/to/sysdig-to-sarif.py report.json > results.sarif

      - name: Upload scan results to GitHub Security tab
        uses: github/codeql-action/upload-sarif@v2
        if: always()
        with:
          sarif_file: results.sarif
```