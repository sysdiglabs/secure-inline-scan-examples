#!/bin/bash
#SECURE_API_TOKEN environment variable must be defined

IMAGE=$1

### Use this block to get JSON output in output.json, as well as "human readable output"
### Begin - no human readable output execution ###
docker run -v /var/run/docker.sock:/var/run/docker.sock -e SYSDIG_API_TOKEN=$SECURE_API_TOKEN quay.io/sysdig/secure-inline-scan:2 $IMAGE --format=JSON > output.json
### End - no human readable output execution ###

### Use this block to get JSON output in output.json, as well as "human readable output" in stdout
### Begin - add human readable output execution ###
# CONTAINER_ID=$(docker run -d --entrypoint /bin/cat -ti -v /var/run/docker.sock:/var/run/docker.sock -e SYSDIG_API_TOKEN=$SECURE_API_TOKEN quay.io/sysdig/secure-inline-scan:2)
# docker exec $CONTAINER_ID mkdir -p /tmp/sysdig-inline-scan/logs/
# docker exec $CONTAINER_ID touch /tmp/sysdig-inline-scan/logs/info.log
# docker exec $CONTAINER_ID tail -f /tmp/sysdig-inline-scan/logs/info.log &
# docker exec $CONTAINER_ID /sysdig-inline-scan.sh $IMAGE --format=JSON > output.json
# exit_status=$?
# sleep 1
# docker stop $CONTAINER_ID -t 0 > /dev/null && docker rm $CONTAINER_ID > /dev/null
### End - add human readable output execution ###

# Check exit status. 0 or 1 is ok to continue (pass or fail policy). Otherwise, report error

exit_status=$?
if [ $exit_status -gt 1 ]; then
    cat output.json
    exit $exit_status
fi

echo "Scan finished. Generating reports"


# Extract vuln report from the output to vulns.json
jq -r '.vulnsReport' output.json > vulns.json

# Create CSV report using mustache
cat <<EOF | docker run -v $(pwd)/vulns.json:/vulns.json --rm -i toolbelt/mustache /vulns.json - > vulns.csv
sep=;
Vuln;Severity;Package;Package_Type;Fix;Url\n" 
{{#vulnerabilities}}
{{vuln}};{{severity}};{{package}};{{package_type}};{{fix}};{{url}}
{{/vulnerabilities}}
EOF

echo "vulns.csv generated"

# Create HTML report using mustache
cat <<EOF | docker run -v $(pwd)/vulns.json:/vulns.json --rm -i toolbelt/mustache /vulns.json - > vulns.html
<html>
<head>
    <title>Vuln report</title>
    <style>
        body {
            font-family: Arial, Helvetica, sans-serif;
        }
        table {
            border-collapse: collapse;
        }
        td, th {
            border: 1px solid black;
            padding: 2px;
        }
    </style>
</head>
<body>
    <table>
        <thead>
            <tr>
                <th>Vuln</th>
                <th>Severity</th>
                <th>Package</th>
                <th>Package_Type</th>
                <th>Fix</th>
                <th>Url</th>
            </tr>
        </thead>
        <tbody>
            {{#vulnerabilities}}
            <tr>
                <td>{{vuln}}</td>
                <td>{{severity}}</td>
                <td>{{package}}</td>
                <td>{{package_type}}</td>
                <td>{{fix}}</td>
                <td><a href="{{url}}">{{url}}</a></td>
            </tr>
            {{/vulnerabilities}}
        </tbody>
    </table>
</body>
</html>
EOF

echo "vulns.html generated"
