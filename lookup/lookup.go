package lookup

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type IPInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

func Scan(urlString string) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
	}

	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		log.Fatal(err)
	}
	ip := ips[0].String()
	url := "https://ipinfo.io/" + ip + "/json"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	var ipInfo IPInfo
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	if err != nil {
		fmt.Println("JSON decoding error:", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})

	table.Append([]string{color.RedString("IP"), color.BlueString(ipInfo.IP)})
	table.Append([]string{color.RedString("Hostname"), color.BlueString(ipInfo.Hostname)})
	table.Append([]string{color.RedString("City"), color.BlueString(ipInfo.City)})
	table.Append([]string{color.RedString("Region"), color.BlueString(ipInfo.Region)})
	table.Append([]string{color.RedString("Country"), color.BlueString(ipInfo.Country)})
	table.Append([]string{color.RedString("Loc"), color.BlueString(ipInfo.Loc)})
	table.Append([]string{color.RedString("Org"), color.BlueString(ipInfo.Org)})
	table.Append([]string{color.RedString("Postal"), color.BlueString(ipInfo.Postal)})
	table.Append([]string{color.RedString("Timezone"), color.BlueString(ipInfo.Timezone)})
	table.Append([]string{color.RedString("Readme"), color.BlueString(ipInfo.Readme)})

	table.Render()
}
