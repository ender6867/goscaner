package vuln

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type WPVersionInfo struct {
	ReleaseDate     string          `json:"release_date"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type Vulnerability struct {
	CVE         string  `json:"cve"`
	CVSS        CVSS    `json:"cvss"`
	Severity    string  `json:"severity"`
	FixedIn     string  `json:"fixed_in"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	References  struct {
		URL []string `json:"url"`
		CVE []string `json:"cve"`
	} `json:"references"`
	Verified bool `json:"verified"`
}

type CVSS struct {
	Score    string `json:"score"`
	Vector   string `json:"vector"`
	Severity string `json:"severity"`
}

func NormalizeVersion(version string) string {
	return strings.ReplaceAll(version, ".", "")
}

func Scan(version string, apiKey string) {
	normalizedVersion := NormalizeVersion(version)
	apiURL := "https://wpscan.com/api/v3/wordpresses/" + normalizedVersion
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Request creation error:", err)
		return
	}

	req.Header.Set("Authorization", "Token token="+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	var data map[string]WPVersionInfo

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("JSON decoding error:", err)
		return
	}

	vulnerabilities := data[version].Vulnerabilities

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Fixed In", "Title", "Description", "CVE References"})

	for _, vuln := range vulnerabilities {
		description := "N/A"
		if vuln.Description != nil {
			description = *vuln.Description
		}

		cveReferences := "N/A"
		if len(vuln.References.CVE) > 0 {
			cveReferences = strings.Join(vuln.References.CVE, ", ")
		}

		table.Append([]string{
			color.GreenString(vuln.FixedIn),
			vuln.Title,
			color.GreenString(description),
			color.GreenString(cveReferences),
		})
	}

	if len(vulnerabilities) > 0 {
		table.Render()
	} else {
		fmt.Println("No vulnerabilities found.")
	}

}
