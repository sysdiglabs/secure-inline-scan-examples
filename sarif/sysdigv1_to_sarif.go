package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

// --- Input Report Structures (based on new report.ts) ---

type Report struct {
	Info    Info    `json:"info"`
	Scanner Scanner `json:"scanner"`
	Result  Result  `json:"result"`
}

type Info struct {
	ScanTime     string `json:"scanTime"`
	ScanDuration string `json:"scanDuration"`
	ResultUrl    string `json:"resultUrl"`
	ResultId     string `json:"resultId"`
}

type Scanner struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Result struct {
	AssetType       string                   `json:"assetType"`
	Layers          map[string]Layer         `json:"layers"`
	Metadata        Metadata                 `json:"metadata"`
	Packages        map[string]Package       `json:"packages"`
	Vulnerabilities map[string]Vulnerability `json:"vulnerabilities"`
}

type Layer struct {
	Command string `json:"command"`
	Digest  string `json:"digest,omitempty"`
	Index   int    `json:"index"`
	Size    int    `json:"size,omitempty"`
}

type Metadata struct {
	Architecture string            `json:"architecture"`
	Autor        string            `json:"autor"`
	BaseOs       string            `json:"baseOs"`
	CreatedAt    string            `json:"createdAt"`
	Digest       string            `json:"digest"`
	ImageId      string            `json:"imageId"`
	Labels       map[string]string `json:"labels,omitempty"`
	Os           string            `json:"os"`
	PullString   string            `json:"pullString"`
	Size         int               `json:"size"`
}

type Package struct {
	IsRemoved           bool     `json:"isRemoved"`
	IsRunning           bool     `json:"isRunning"`
	LayerRef            string   `json:"layerRef"`
	Name                string   `json:"name"`
	Path                string   `json:"path"`
	SuggestedFix        string   `json:"suggestedFix,omitempty"`
	Type                string   `json:"type"`
	Version             string   `json:"version"`
	VulnerabilitiesRefs []string `json:"vulnerabilitiesRefs"`
}

type Vulnerability struct {
	CvssScore      CvssScore `json:"cvssScore"`
	DisclosureDate string    `json:"disclosureDate"`
	Exploitable    bool      `json:"exploitable"`
	FixVersion     string    `json:"fixVersion,omitempty"`
	Name           string    `json:"name"`
	PackageRef     string    `json:"packageRef"`
	RiskAcceptRefs []string  `json:"riskAcceptRefs,omitempty"`
	Severity       string    `json:"severity"`
}

type CvssScore struct {
	Version string  `json:"version"`
	Score   float64 `json:"score"`
	Vector  string  `json:"vector"`
}

// --- SARIF Structures (mostly unchanged) ---

type SARIF struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SARIFRun `json:"runs"`
}
type SARIFRun struct {
	Tool             SARIFTool              `json:"tool"`
	LogicalLocations []LogicalLocation      `json:"logicalLocations"`
	Results          []SARIFResult          `json:"results"`
	ColumnKind       string                 `json:"columnKind"`
	Properties       map[string]interface{} `json:"properties"`
}
type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}
type SARIFDriver struct {
	Name                  string      `json:"name"`
	FullName              string      `json:"fullName"`
	InformationUri        string      `json:"informationUri"`
	Version               string      `json:"version"`
	SemanticVersion       string      `json:"semanticVersion"`
	DottedQuadFileVersion string      `json:"dottedQuadFileVersion"`
	Rules                 []SARIFRule `json:"rules"`
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

// --- Filtering & Helpers ---

type FilterOptions struct {
	MinSeverity     string
	PackageTypes    []string
	NotPackageTypes []string
	ExcludeAccepted bool
}

var severityOrder = []string{"negligible", "low", "medium", "high", "critical"}
var severityNames = []string{"critical", "high", "medium", "low", "negligible"}

func getSeverityIndex(severity string) int {
	lowerSev := strings.ToLower(severity)
	for i, s := range severityOrder {
		if s == lowerSev {
			return i
		}
	}
	return -1
}

func isSeverityGte(a, b string) bool {
	return getSeverityIndex(a) >= getSeverityIndex(b)
}

func numericPriorityForSeverity(severity string) int {
	lowerSev := strings.ToLower(severity)
	for i, s := range severityNames {
		if s == lowerSev {
			return i
		}
	}
	return 5 // Default to lowest priority
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
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

func filterPackages(pkgs map[string]Package, vulns map[string]Vulnerability, opts FilterOptions) map[string]Package {
	filtered := make(map[string]Package)

PackageLoop:
	for key, pkg := range pkgs {
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
				continue PackageLoop
			}
		}
		for _, notAllowed := range opts.NotPackageTypes {
			if strings.ToLower(notAllowed) == ptype {
				continue PackageLoop
			}
		}

		var newVulnRefs []string
		for _, vulnRef := range pkg.VulnerabilitiesRefs {
			vuln, ok := vulns[vulnRef]
			if !ok {
				continue
			}
			if opts.MinSeverity != "" && !isSeverityGte(vuln.Severity, opts.MinSeverity) {
				continue
			}
			if opts.ExcludeAccepted && len(vuln.RiskAcceptRefs) > 0 {
				continue
			}
			newVulnRefs = append(newVulnRefs, vulnRef)
		}

		if len(newVulnRefs) > 0 {
			newPkg := pkg
			newPkg.VulnerabilitiesRefs = newVulnRefs
			filtered[key] = newPkg
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

func getSARIFVulnFullDescription(pkg Package, vuln Vulnerability) string {
	fix := "No fix available"
	if pkg.SuggestedFix != "" {
		fix = pkg.SuggestedFix
	}
	return fmt.Sprintf("%s\nSeverity: %s\nPackage: %s\nType: %s\nFix: %s\nURL: https://nvd.nist.gov/vuln/detail/%s",
		vuln.Name, vuln.Severity, pkg.Name, pkg.Type, fix, vuln.Name)
}

func getSARIFPkgHelp(pkg Package, vulns map[string]Vulnerability) SARIFHelp {
	var textBuilder, markdownBuilder strings.Builder

	markdownBuilder.WriteString("| Vulnerability | Severity | CVSS Score | CVSS Version | CVSS Vector | Exploitable |\n")
	markdownBuilder.WriteString("| -------- | ------- | ---------- | ------------ | ----------- | ----------- |\n")

	// Sort vulnerabilities for consistent output
	sort.Strings(pkg.VulnerabilitiesRefs)

	for _, vulnRef := range pkg.VulnerabilitiesRefs {
		vuln := vulns[vulnRef]
		fix := "No fix available"
		if pkg.SuggestedFix != "" {
			fix = pkg.SuggestedFix
		}

		textBuilder.WriteString(fmt.Sprintf("Vulnerability %s\n  Severity: %s\n  Package: %s\n  CVSS Score: %.1f\n  CVSS Version: %s\n  CVSS Vector: %s\n  Version: %s\n  Fix Version: %s\n  Exploitable: %t\n  Type: %s\n  Location: %s\n  URL: https://nvd.nist.gov/vuln/detail/%s\n\n\n",
			vuln.Name, vuln.Severity, pkg.Name, vuln.CvssScore.Score, vuln.CvssScore.Version, vuln.CvssScore.Vector, pkg.Version, fix, vuln.Exploitable, pkg.Type, pkg.Path, vuln.Name,
		))

		markdownBuilder.WriteString(fmt.Sprintf("| %s | %s | %.1f | %s | %s | %t |\n",
			vuln.Name, vuln.Severity, vuln.CvssScore.Score, vuln.CvssScore.Version, vuln.CvssScore.Vector, vuln.Exploitable,
		))
	}

	return SARIFHelp{
		Text:     textBuilder.String(),
		Markdown: markdownBuilder.String(),
	}
}

func sanitizeImageName(imageName string) string {
	re := regexp.MustCompile(`[/:_]`)
	return re.ReplaceAllString(imageName, "-")
}

// --- Main Conversion Logic ---

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
			"layersCount":  len(data.Result.Layers),
			"resultUrl":    data.Info.ResultUrl,
			"resultId":     data.Info.ResultId,
		},
	}
	run.Tool.Driver.Name = "sysdig-cli-scanner"
	run.Tool.Driver.FullName = "Sysdig Vulnerability CLI Scanner"
	run.Tool.Driver.InformationUri = "https://docs.sysdig.com/en/docs/installation/sysdig-secure/install-vulnerability-cli-scanner"
	run.Tool.Driver.Version = "1.0.0"         // Placeholder version
	run.Tool.Driver.SemanticVersion = "1.0.0" // Placeholder version
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

	// Sort packages by name for consistent output
	pkgKeys := make([]string, 0, len(data.Result.Packages))
	for k := range data.Result.Packages {
		pkgKeys = append(pkgKeys, k)
	}
	sort.Strings(pkgKeys)

	for _, key := range pkgKeys {
		pkg := data.Result.Packages[key]
		if len(pkg.VulnerabilitiesRefs) == 0 {
			continue
		}

		var fullDescBuilder strings.Builder
		score := 0.0
		severityLevel := ""
		minSeverityNum := 5

		for _, vulnRef := range pkg.VulnerabilitiesRefs {
			vuln, ok := data.Result.Vulnerabilities[vulnRef]
			if !ok {
				continue
			}
			fullDescBuilder.WriteString(getSARIFVulnFullDescription(pkg, vuln) + "\n\n")

			if vuln.CvssScore.Score > score {
				score = vuln.CvssScore.Score
			}

			sevNum := numericPriorityForSeverity(vuln.Severity)
			if sevNum < minSeverityNum {
				severityLevel = strings.ToLower(vuln.Severity)
				minSeverityNum = sevNum
			}
		}

		rules = append(rules, SARIFRule{
			Id:               pkg.Name,
			Name:             pkg.Name,
			ShortDescription: SARIFText{Text: fmt.Sprintf("Vulnerable package: %s", pkg.Name)},
			FullDescription:  SARIFText{Text: strings.TrimSpace(fullDescBuilder.String())},
			HelpUri:          "", // No direct equivalent in new model
			Help:             getSARIFPkgHelp(pkg, data.Result.Vulnerabilities),
			Properties: SARIFProperties{
				Precision:        "very-high",
				SecuritySeverity: fmt.Sprintf("%.1f", score),
				Tags:             []string{"vulnerability", "security", severityLevel},
			},
		})

		results = append(results, SARIFResult{
			RuleId: pkg.Name,
			Level:  checkLevel(severityLevel),
			Message: SARIFText{
				Text: fmt.Sprintf("Vulnerabilities found in package %s", pkg.Name),
			},
			Locations: []SARIFLocation{{
				PhysicalLocation: SARIFPhysicalLocation{
					ArtifactLocation: SARIFArtifactLocation{
						Uri:       "file:///" + sanitizeImageName(data.Result.Metadata.PullString),
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
	seenRules := make(map[string]bool)

	// Sort packages and vulnerabilities for consistent output
	pkgKeys := make([]string, 0, len(data.Result.Packages))
	for k := range data.Result.Packages {
		pkgKeys = append(pkgKeys, k)
	}
	sort.Strings(pkgKeys)

	for _, key := range pkgKeys {
		pkg := data.Result.Packages[key]
		sort.Strings(pkg.VulnerabilitiesRefs)

		for _, vulnRef := range pkg.VulnerabilitiesRefs {
			vuln, ok := data.Result.Vulnerabilities[vulnRef]
			if !ok {
				continue
			}

			if !seenRules[vuln.Name] {
				seenRules[vuln.Name] = true
				rules = append(rules, SARIFRule{
					Id:               vuln.Name,
					Name:             pkg.Type,
					ShortDescription: SARIFText{Text: fmt.Sprintf("%s Severity: %s Package: %s", vuln.Name, vuln.Severity, pkg.Name)},
					FullDescription:  SARIFText{Text: getSARIFVulnFullDescription(pkg, vuln)},
					HelpUri:          "https://nvd.nist.gov/vuln/detail/" + vuln.Name,
					Help:             SARIFHelp{Text: "See NVD", Markdown: ""}, // Simplified help for individual vulns
					Properties: SARIFProperties{
						Precision:        "very-high",
						SecuritySeverity: fmt.Sprintf("%.1f", vuln.CvssScore.Score),
						Tags:             []string{"vulnerability", "security", strings.ToLower(vuln.Severity)},
					},
				})
			}

			results = append(results, SARIFResult{
				RuleId: vuln.Name,
				Level:  checkLevel(vuln.Severity),
				Message: SARIFText{
					Text: fmt.Sprintf("Vulnerability %s in package %s", vuln.Name, pkg.Name),
				},
				Locations: []SARIFLocation{{
					PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{
							Uri:       "file:///" + sanitizeImageName(data.Result.Metadata.PullString),
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

// --- Main execution ---

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
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", inputFile, err)
		os.Exit(1)
	}

	var data Report
	if err := json.Unmarshal(raw, &data); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling JSON: %v\n", err)
		os.Exit(1)
	}

	opts := FilterOptions{
		MinSeverity:     *minSeverity,
		PackageTypes:    splitAndTrim(*packageTypes),
		NotPackageTypes: splitAndTrim(*notPackageTypes),
		ExcludeAccepted: *excludeAccepted,
	}
	data.Result.Packages = filterPackages(data.Result.Packages, data.Result.Vulnerabilities, opts)

	sarif := vulnerabilities2SARIF(data, *groupByPackage)
	out, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling SARIF: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(outputFile, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing SARIF file %s: %v\n", outputFile, err)
		os.Exit(1)
	}
	fmt.Println("SARIF written to", outputFile)
}
