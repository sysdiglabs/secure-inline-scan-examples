#!/usr/bin/env python3

import argparse
import json
import logging

# Setup logger
LOG = logging.getLogger(__name__)

# define a Handler which writes INFO messages or higher to the sys.stderr
console = logging.StreamHandler()
# set a format which is simpler for console use
console.setLevel(logging.INFO)
# tell the handler to use this format
console.setFormatter(logging.Formatter('%(message)s'))

# add the handler to the main logger
LOG.addHandler(console)


LEVELS = {
    "error": ["High","Critical"],
    "warning": ["Medium"],
    "note": ["Negligible","Low"]
}

def generate_report(data):
    results = []
    rules = []
    ruleIds = []
    resultUrl = data['info']['resultUrl']
    baseUrl = resultUrl[:resultUrl.rfind('/')]
    for package in data['result']['packages']:
        # Continue if package has no vulnerabilities
        if 'vulns' not in package.keys():
            LOG.info(f"Package: {package['name']} has no vulnerabilities...skipping...")
            continue
        
        for vuln in package['vulns']:
            if vuln['name'] not in ruleIds:
                ruleIds.append(vuln['name'])
                rule = {
                    "id": f"{vuln['name']}",
                    "name": f"{package['type']}",
                    "shortDescription": {
                        "text": f"{vuln['name']} - {package['name']}@{package['version']}"
                    },
                    "fullDescription": {
                        "text": f"{vuln['name']} - {package['name']}@{package['version']}"
                    },
                    "defaultConfiguration": {
                        "level": f"{check_level(vuln['severity']['value'])}"
                    },
                    "helpUri": f"https://nvd.nist.gov/vuln/detail/{vuln['name']}",
                    "help": {
                        "text": f"Vulnerability {vuln['name']}\nPackage: {package['name']}\nSeverity: {vuln['severity']['value']}\nCVSS Score: {vuln['cvssScore']['value']['score']}\nCVSS Version: {vuln['cvssScore']['value']['version']}\nCVSS Vector: {vuln['cvssScore']['value']['vector']}\nFixed Version: {(vuln['fixedInVersion'] if 'fixedInVersion' in vuln else '')}\nExploitable: {vuln['exploitable']}\nLink: [{vuln['name']}](https://nvd.nist.gov/vuln/detail/{vuln['name']})",
                        "markdown": f"**Vulnerability {vuln['name']}**\n| Package | Severity| CVSS Score | CVSS Version | CVSS Vector | Fixed Version | Exploitable | Link |\n| --- | --- | --- | --- | --- | --- | --- | --- |\n|{package['name']}|{vuln['severity']['value']}|{vuln['cvssScore']['value']['score']}|{vuln['cvssScore']['value']['version']}|{vuln['cvssScore']['value']['vector']}|{(vuln['fixedInVersion'] if 'fixedInVersion' in vuln else '')}|{vuln['exploitable']}|[{vuln['name']}](https://nvd.nist.gov/vuln/detail/{vuln['name']})|"
                    },
                    "properties": {
                        "precision": "very-high",
                        "security-severity": f"{vuln['cvssScore']['value']['score']}",
                        "tags": [
                            "vulnerability",
                            "security",
                            f"{vuln['severity']['value']}"
                        ]
                    }
                }
                rules.append(rule)

            result = {
                "ruleId": f"{vuln['name']}",
                "level": f"{check_level(vuln['severity']['value'])}",
                "message": {
                    "text": f"Full image scan results in Sysdig UI: [{data['result']['metadata']['pullString']} scan result]({data['info']['resultUrl']})\nPackage: [{package['name']}]({baseUrl}/content?filter=freeText+in+(\"{package['name']}\"))\nPackage type: {package['type']}\nInstalled Version: {package['version']}\nPackage path: {package['path']}\nVulnerability: [{vuln['name']}]({baseUrl}/vulnerabilities?filter=freeText+in+(\"{vuln['name']}\"))\nSeverity: {vuln['severity']['value']}\nCVSS Score: {vuln['cvssScore']['value']['score']}\nCVSS Version: {vuln['cvssScore']['value']['version']}\nCVSS Vector: {vuln['cvssScore']['value']['vector']}\nFixed Version: {(vuln['fixedInVersion'] if 'fixedInVersion' in vuln else '')}\nExploitable: {vuln['exploitable']}\nLink to NVD: [{vuln['name']}](https://nvd.nist.gov/vuln/detail/{vuln['name']})"
                },
                "locations": [
                    {
                        "physicalLocation": {
                            "artifactLocation": {
                                "uri": f"{data['result']['metadata']['pullString']}",
                                "uriBaseId": "ROOTPATH"
                            }
                        },
                        "message": {
                            "text": f"{data['result']['metadata']['pullString']} - {package['name']}@{package['version']}"
                        }
                    }
                ]
            }
            results.append(result)

    run = {
        "tool": {
            "driver": {
                "fullName": "Sysdig Vulnerability CLI Scanner",
                "informationUri": "https://docs.sysdig.com/en/docs/installation/sysdig-secure/install-vulnerability-cli-scanner",
                "name": "sysdig-cli-scanner",
                "version": f"{data['scanner']['version']}",
                "rules": rules
            }
        },
        "results": results,
        "columnKind": "utf16CodeUnits",
        "properties": {
            "pullString": f"{data['result']['metadata']['pullString']}",
            "digest": f"{data['result']['metadata']['digest']}",
            "imageId": f"{data['result']['metadata']['imageId']}",
            "architecture": f"{data['result']['metadata']['architecture']}",
            "baseOs": f"{data['result']['metadata']['baseOs']}",
            "os": f"{data['result']['metadata']['os']}",
            "size": f"{data['result']['metadata']['size']}",
            "layersCount": f"{data['result']['metadata']['layersCount']}",
            "resultUrl": f"{data['info']['resultUrl']}",
            "resultId": f"{data['info']['resultId']}",
        }
    }

    report = {
        "version": "2.1.0",
        "$schema": "https://json.schemastore.org/sarif-2.1.0.json",
        "runs": [run]
    }

    return report

def check_level(severity):
    for key in LEVELS:
        if severity in LEVELS[key]:
            return key

def main():
    parser = argparse.ArgumentParser(description="Convert Sysdig report to SARIF format.")
    parser.add_argument("filename", type=str, help="Sysdig report file in json format.")
    parser.add_argument("--log-level", dest="logLevel", type=str, choices=['INFO', 'DEBUG'], help="Set log level. If DEBUG level set logs are stored in report.log file in the same folder where script is executed.", default="INFO", required=False)
    args = parser.parse_args()
    filename = args.filename
    logLevel = args.logLevel

    if logLevel == "DEBUG":
        logging.basicConfig(
            level=logging.DEBUG,
            format="%(asctime)s.%(msecs)03d %(levelname)s - %(funcName)s: %(message)s",
            datefmt="%Y-%m-%d %H:%M:%S",
            filename="report.log",
            filemode = 'w',
        )
    else:
        LOG.setLevel(logging.INFO)

    try:
      LOG.info(f"Loading {filename} JSON file.")
      with open(filename) as json_file:
        data = json.load(json_file)

    except:
      LOG.info(f"Error: {filename}: Invalid JSON file!")
      exit(parser.print_help())

    # Simple validation of JSON file format
    try:
        'metadata' in data
        'vulnerabilties' in data
        'packages' in data
        'policies' in data
        'info' in data
    except:
      LOG.info(f"Error: {filename}: JSON file is not from sysdig-cli-scanner!")
      exit(parser.print_help())

    report = generate_report(data=data)
    print(json.dumps(report))
    
if __name__ == '__main__':
    main()
