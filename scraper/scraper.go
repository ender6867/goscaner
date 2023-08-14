package scraper

import (
	"fmt"
	"goscaner/vuln"
	"log"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
)

func Scrape(urlString string) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
	}

	host := parsedURL.Hostname()

	ipAddress, err := net.LookupIP(host)
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector()

	var wordpressVersion string
	var themeNames []string
	var pluginNames []string

	var mu sync.Mutex

	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, "/wp-content/themes/") {
			parts := strings.Split(href, "/wp-content/themes/")
			if len(parts) > 1 {
				themeName := strings.Split(parts[1], "/")[0]
				mu.Lock()
				if !contains(themeNames, themeName) {
					themeNames = append(themeNames, themeName)
				}
				mu.Unlock()
			}
		} else if strings.Contains(href, "/wp-content/plugins/") {
			parts := strings.Split(href, "/wp-content/plugins/")
			if len(parts) > 1 {
				pluginName := strings.Split(parts[1], "/")[0]
				mu.Lock()
				if !contains(pluginNames, pluginName) {
					pluginNames = append(pluginNames, pluginName)
				}
				mu.Unlock()
			}
		}
	})

	c.OnHTML("meta[name='generator']", func(e *colly.HTMLElement) {
		content := e.Attr("content")

		re := regexp.MustCompile(`WordPress (\d+\.\d+(\.\d+)?)`)
		match := re.FindStringSubmatch(content)
		if len(match) >= 2 {
			wordpressVersion = match[1]
		}
	})

	err = c.Visit(urlString)
	if err != nil {
		log.Fatal(err)
	}

	printResults(ipAddress, themeNames, pluginNames, wordpressVersion)
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func printResults(ips []net.IP, themes []string, plugins []string, version string) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	table := tablewriter.NewWriter(color.Output)
	table.SetHeader([]string{"Type", "Value"})

	ipStrings := make([]string, len(ips))
	for i, ip := range ips {
		ipStrings[i] = ip.String()
	}
	ipValue := strings.Join(ipStrings, ", ")
	table.Append([]string{color.RedString("IP_Address "), color.RedString(ipValue)})

	if len(themes) > 0 {
		themesString := strings.Join(themes, ", ")
		table.Append([]string{color.CyanString("Themes"), color.CyanString(themesString)})
	} else {
		table.Append([]string{"Themes", color.YellowString("No themes found")})
	}

	if len(plugins) > 0 {
		pluginsString := strings.Join(plugins, ", ")
		table.Append([]string{color.MagentaString("Plugins"), color.MagentaString(pluginsString)})
	} else {
		table.Append([]string{"Plugins", color.YellowString("No plugins found")})
	}

	table.Append([]string{color.BlueString("WordPress_Version"), color.BlueString(version)})

	table.Render()
	apiKey := os.Getenv("WPSCAN_API_KEY")
	if version != "" {
		vuln.Scan(version, apiKey)
	}

}
