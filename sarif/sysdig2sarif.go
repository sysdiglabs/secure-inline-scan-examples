package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ----------- Sysdig structs (simplified; adjust fields to match your actual JSON) -----------

type Report struct {
	Info   Info   `json:"info"`
	Result Result `json:"result"`
}

type Info struct {
	ResultUrl string `json:"resultUrl"`
	ResultId  string `json:"resultId"`
}

type Result struct {
	Metadata Metadata  `json:"metadata"`
	Packages []Package `json:"packages"`
}

type Metadata struct {
	PullString   string `json:"pullString"`
	Digest       string `json:"digest"`
	ImageId      string `json:"imageId"`
	Architecture string `json:"architecture"`
	BaseOs       string `json:"baseOs"`
	Os           string `json:"os"`
	Size         int    `json:"size"`
	LayersCount  int    `json:"layersCount"`
}

type Package struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Type         string `json:"type"`
	Path         string `json:"path"`
	SuggestedFix string `json:"suggestedFix"`
	Vulns        []Vuln `json:"vulns"`
}

type Vuln struct {
	Name           string    `json:"name"`
	Severity       Severity  `json:"severity"`
	CvssScore      CvssScore `json:"cvssScore"`
	Exploitable    bool      `json:"exploitable"`
	FixedInVersion string    `json:"fixedInVersion"`
}

type Severity struct {
	Value string `json:"value"`
}

type CvssScore struct {
	Value CvssValue `json:"value"`
}

type CvssValue struct {
	Score   float64 `json:"score"`
	Version string  `json:"version"`
	Vector  string  `json:"vector"`
}

// ----------- SARIF structs (essential subset) -----------

type SARIF struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool struct {
		Driver struct {
			Name                  string      `json:"name"`
			FullName              string      `json:"fullName"`
			InformationUri        string      `json:"informationUri"`
			Version               string      `json:"version"`
			SemanticVersion       string      `json:"semanticVersion"`
			DottedQuadFileVersion string      `json:"dottedQuadFileVersion"`
			Rules                 []SARIFRule `json:"rules"`
		} `json:"driver"`
	} `json:"tool"`
	LogicalLocations []LogicalLocation      `json:"logicalLocations"`
	Results          []SARIFResult          `json:"results"`
	ColumnKind       string                 `json:"columnKind"`
	Properties       map[string]interface{} `json:"properties"`
}

type LogicalLocation struct {
	Name               string `json:"name"`
	FullyQualifiedName string `json:"fullyQualifiedName"`
	Kind               string `json:"kind"`
}

type SARIFRule struct {
	Id               string          `json:"id"`
	Name             string          `json:"name"`
	ShortDescription SARIFText       `json:"shortDescription"`
	FullDescription  SARIFText       `json:"fullDescription"`
	HelpUri          string          `json:"helpUri"`
	Help             SARIFHelp       `json:"help"`
	Properties       SARIFProperties `json:"properties"`
}

type SARIFText struct {
	Text string `json:"text"`
}

type SARIFHelp struct {
	Text     string `json:"text"`
	Markdown string `json:"markdown"`
}

type SARIFProperties struct {
	Precision        string   `json:"precision"`
	SecuritySeverity string   `json:"security-severity"`
	Tags             []string `json:"tags"`
}

type SARIFResult struct {
	RuleId    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   SARIFText       `json:"message"`
	Locations []SARIFLocation `json:"locations"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
	Message          SARIFText             `json:"message"`
}

type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
}

type SARIFArtifactLocation struct {
	Uri       string `json:"uri"`
	UriBaseId string `json:"uriBaseId"`
}

// ---------- Helpers ----------

// Converts Sysdig severity string to SARIF level
func checkLevel(sev string) string {
	sev = strings.ToLower(sev)
	switch sev {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	case "low", "negligible":
		return "note"
	default:
		return "note"
	}
}

// Formats a vulnerability's full description for SARIF
func getVulnFullDescription(pkg Package, vuln Vuln) string {
	return fmt.Sprintf("%s\nSeverity: %s\nPackage: %s\nType: %s\nFix: %s\nURL: https://nvd.nist.gov/vuln/detail/%s",
		vuln.Name, vuln.Severity.Value, pkg.Name, pkg.Type, pkg.SuggestedFix, vuln.Name)
}

// ---------- Main CLI Entry Point ----------

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sysdig2sarif input.json output.sarif [groupByPackage]")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	groupByPackage := false
	if len(os.Args) >= 4 && os.Args[3] == "true" {
		groupByPackage = true
	}

	raw, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	var data Report
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(err)
	}

	sarif := vulnerabilities2SARIF(data, groupByPackage)
	out, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(outputFile, out, 0644); err != nil {
		panic(err)
	}
	fmt.Println("SARIF written to", outputFile)
}

// ---------- Conversion Functions ----------

// Converts the Sysdig report to a SARIF object
func vulnerabilities2SARIF(data Report, groupByPackage bool) SARIF {
	var rules []SARIFRule
	var results []SARIFResult

	if groupByPackage {
		rules, results = vulnerabilities2SARIFResByPackage(data)
	} else {
		rules, results = vulnerabilities2SARIFRes(data)
	}

	run := SARIFRun{
		LogicalLocations: []LogicalLocation{{
			Name:               "container-image",
			FullyQualifiedName: "container-image",
			Kind:               "namespace",
		}},
		Results:    results,
		ColumnKind: "utf16CodeUnits",
		Properties: map[string]interface{}{
			"pullString":   data.Result.Metadata.PullString,
			"digest":       data.Result.Metadata.Digest,
			"imageId":      data.Result.Metadata.ImageId,
			"architecture": data.Result.Metadata.Architecture,
			"baseOs":       data.Result.Metadata.BaseOs,
			"os":           data.Result.Metadata.Os,
			"size":         data.Result.Metadata.Size,
			"layersCount":  data.Result.Metadata.LayersCount,
			"resultUrl":    data.Info.ResultUrl,
			"resultId":     data.Info.ResultId,
		},
	}
	// Fill in tool/driver info
	run.Tool.Driver.Name = "sysdig-cli-scanner"
	run.Tool.Driver.FullName = "Sysdig Vulnerability CLI Scanner"
	run.Tool.Driver.InformationUri = "https://docs.sysdig.com/en/docs/installation/sysdig-secure/install-vulnerability-cli-scanner"
	run.Tool.Driver.Version = "1.0.0" // Change as appropriate
	run.Tool.Driver.SemanticVersion = "1.0.0"
	run.Tool.Driver.DottedQuadFileVersion = "1.0.0.0"
	run.Tool.Driver.Rules = rules

	return SARIF{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs:    []SARIFRun{run},
	}
}

// SARIF conversion, grouping results by package
func vulnerabilities2SARIFResByPackage(data Report) ([]SARIFRule, []SARIFResult) {
	var rules []SARIFRule
	var results []SARIFResult

	for _, pkg := range data.Result.Packages {
		if len(pkg.Vulns) == 0 {
			continue
		}
		fullDesc := ""
		score := 0.0
		severity := ""
		for _, vuln := range pkg.Vulns {
			fullDesc += getVulnFullDescription(pkg, vuln) + "\n\n"
			if vuln.CvssScore.Value.Score > score {
				score = vuln.CvssScore.Value.Score
			}
			if severity == "" {
				severity = vuln.Severity.Value
			}
		}
		rules = append(rules, SARIFRule{
			Id:               pkg.Name,
			Name:             pkg.Name,
			ShortDescription: SARIFText{Text: fmt.Sprintf("Vulnerable package: %s", pkg.Name)},
			FullDescription:  SARIFText{Text: fullDesc},
			HelpUri:          "",
			Help:             SARIFHelp{Text: "Multiple vulnerabilities", Markdown: ""},
			Properties: SARIFProperties{
				Precision:        "very-high",
				SecuritySeverity: fmt.Sprintf("%v", score),
				Tags:             []string{"vulnerability", "security", severity},
			},
		})

		results = append(results, SARIFResult{
			RuleId: pkg.Name,
			Level:  checkLevel(severity),
			Message: SARIFText{
				Text: fmt.Sprintf("Vulnerabilities found in package %s", pkg.Name),
			},
			Locations: []SARIFLocation{{
				PhysicalLocation: SARIFPhysicalLocation{
					ArtifactLocation: SARIFArtifactLocation{
						Uri:       "file:///" + data.Result.Metadata.PullString,
						UriBaseId: "ROOTPATH",
					},
				},
				Message: SARIFText{
					Text: fmt.Sprintf("%s - %s@%s", data.Result.Metadata.PullString, pkg.Name, pkg.Version),
				},
			}},
		})
	}
	return rules, results
}

// SARIF conversion, result per vulnerability
func vulnerabilities2SARIFRes(data Report) ([]SARIFRule, []SARIFResult) {
	var rules []SARIFRule
	var results []SARIFResult
	seen := map[string]bool{}

	for _, pkg := range data.Result.Packages {
		for _, vuln := range pkg.Vulns {
			if !seen[vuln.Name] {
				seen[vuln.Name] = true
				rules = append(rules, SARIFRule{
					Id:               vuln.Name,
					Name:             pkg.Type,
					ShortDescription: SARIFText{Text: fmt.Sprintf("%s Severity: %s Package: %s", vuln.Name, vuln.Severity.Value, pkg.Name)},
					FullDescription:  SARIFText{Text: getVulnFullDescription(pkg, vuln)},
					HelpUri:          "https://nvd.nist.gov/vuln/detail/" + vuln.Name,
					Help:             SARIFHelp{Text: "See NVD", Markdown: ""},
					Properties: SARIFProperties{
						Precision:        "very-high",
						SecuritySeverity: fmt.Sprintf("%v", vuln.CvssScore.Value.Score),
						Tags:             []string{"vulnerability", "security", vuln.Severity.Value},
					},
				})
			}
			results = append(results, SARIFResult{
				RuleId: vuln.Name,
				Level:  checkLevel(vuln.Severity.Value),
				Message: SARIFText{
					Text: fmt.Sprintf("Vulnerability %s in package %s", vuln.Name, pkg.Name),
				},
				Locations: []SARIFLocation{{
					PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{
							Uri:       data.Result.Metadata.PullString,
							UriBaseId: "ROOTPATH",
						},
					},
					Message: SARIFText{
						Text: fmt.Sprintf("%s - %s@%s", data.Result.Metadata.PullString, pkg.Name, pkg.Version),
					},
				}},
			})
		}
	}
	return rules, results
}
