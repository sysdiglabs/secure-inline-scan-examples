package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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
	Name           string        `json:"name"`
	Severity       Severity      `json:"severity"`
	CvssScore      CvssScore     `json:"cvssScore"`
	Exploitable    bool          `json:"exploitable"`
	FixedInVersion string        `json:"fixedInVersion"`
	AcceptedRisks  []interface{} `json:"acceptedRisks"`
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

// SARIF structures
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

// Filtering options
type FilterOptions struct {
	MinSeverity     string
	PackageTypes    []string
	NotPackageTypes []string
	ExcludeAccepted bool
}

// Severity order
var severityOrder = []string{"Negligible", "Low", "Medium", "High", "Critical"}

func isSeverityGte(a, b string) bool {
	var ai, bi int = -1, -1
	for i, v := range severityOrder {
		if strings.EqualFold(v, a) {
			ai = i
		}
		if strings.EqualFold(v, b) {
			bi = i
		}
	}
	return ai >= bi && ai >= 0 && bi >= 0
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	var out []string
	for _, v := range parts {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func filterPackages(pkgs []Package, opts FilterOptions) []Package {
	var filtered []Package
TypeLoop:
	for _, pkg := range pkgs {
		ptype := strings.ToLower(pkg.Type)
		if len(opts.PackageTypes) > 0 {
			found := false
			for _, allowed := range opts.PackageTypes {
				if strings.ToLower(allowed) == ptype {
					found = true
					break
				}
			}
			if !found {
				continue TypeLoop
			}
		}
		for _, notAllowed := range opts.NotPackageTypes {
			if strings.ToLower(notAllowed) == ptype {
				continue TypeLoop
			}
		}
		newVulns := []Vuln{}
		for _, vuln := range pkg.Vulns {
			if opts.MinSeverity != "" && !isSeverityGte(vuln.Severity.Value, opts.MinSeverity) {
				continue
			}
			if opts.ExcludeAccepted && len(vuln.AcceptedRisks) > 0 {
				fmt.Printf("Accepted risks: %d\n", len(vuln.AcceptedRisks))
				continue
			}
			newVulns = append(newVulns, vuln)
		}
		if len(newVulns) > 0 {
			pkg.Vulns = newVulns
			filtered = append(filtered, pkg)
		}
	}
	return filtered
}

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

func getVulnFullDescription(pkg Package, vuln Vuln) string {
	return fmt.Sprintf("%s\nSeverity: %s\nPackage: %s\nType: %s\nFix: %s\nURL: https://nvd.nist.gov/vuln/detail/%s",
		vuln.Name, vuln.Severity.Value, pkg.Name, pkg.Type, pkg.SuggestedFix, vuln.Name)
}

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
	run.Tool.Driver.Name = "sysdig-cli-scanner"
	run.Tool.Driver.FullName = "Sysdig Vulnerability CLI Scanner"
	run.Tool.Driver.InformationUri = "https://docs.sysdig.com/en/docs/installation/sysdig-secure/install-vulnerability-cli-scanner"
	run.Tool.Driver.Version = "1.0.0"
	run.Tool.Driver.SemanticVersion = "1.0.0"
	run.Tool.Driver.DottedQuadFileVersion = "1.0.0.0"
	run.Tool.Driver.Rules = rules

	return SARIF{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs:    []SARIFRun{run},
	}
}

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

func main() {
	minSeverity := flag.String("min-severity", "", "Minimum severity (e.g., High)")
	packageTypes := flag.String("type", "", "Package types (comma-separated, e.g., java,javascript)")
	notPackageTypes := flag.String("not-type", "", "Exclude package types (comma-separated)")
	excludeAccepted := flag.Bool("exclude-accepted", false, "Exclude vulnerabilities with accepted risks")
	groupByPackage := flag.Bool("group-by-package", false, "Group by package")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: sysdig2sarif [flags] input.json output.sarif")
		os.Exit(1)
	}
	inputFile := flag.Arg(0)
	outputFile := flag.Arg(1)

	raw, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	var data Report
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(err)
	}

	opts := FilterOptions{
		MinSeverity:     *minSeverity,
		PackageTypes:    splitAndTrim(*packageTypes),
		NotPackageTypes: splitAndTrim(*notPackageTypes),
		ExcludeAccepted: *excludeAccepted,
	}
	data.Result.Packages = filterPackages(data.Result.Packages, opts)

	sarif := vulnerabilities2SARIF(data, *groupByPackage)
	out, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(outputFile, out, 0644); err != nil {
		panic(err)
	}
	fmt.Println("SARIF written to", outputFile)
}
