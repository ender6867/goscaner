package dirscanner

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func Scan(urlString string, wordlist []string) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
	}

	var foundURLs200 [][]string
	var foundURLs302 [][]string
	var foundURLsOther [][]string
	var mu sync.Mutex

	client := http.Client{
		Timeout: time.Second * 10,
	}

	totalRequests := len(wordlist)
	completedRequests := 0
	percentCompleted := 0
	progressInterval := totalRequests / 100

	var wg sync.WaitGroup

	for _, word := range wordlist {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()

			targetURL := parsedURL.ResolveReference(&url.URL{Path: word}).String()
			resp, err := client.Get(targetURL)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			statusText := ""
			if resp.StatusCode == 200 {
				statusText = "200_OK"
				mu.Lock()
				foundURLs200 = append(foundURLs200, []string{color.GreenString(targetURL), color.GreenString(statusText)})
				mu.Unlock()
			} else if resp.StatusCode == 302 {
				statusText = "302_Found"
				mu.Lock()
				foundURLs302 = append(foundURLs302, []string{color.YellowString(targetURL), color.YellowString(statusText)})
				mu.Unlock()

			} else if resp.StatusCode < 400 || resp.StatusCode >= 500 {
				statusText = fmt.Sprintf("%d", resp.StatusCode)
				mu.Lock()
				foundURLsOther = append(foundURLsOther, []string{color.RedString(targetURL), color.RedString(statusText)})
				mu.Unlock()
			}

			completedRequests++

			if completedRequests%progressInterval == 0 {
				percentCompleted++
				printProgress(percentCompleted)
			}
		}(word)
	}

	wg.Wait()

	printResults(foundURLs200)
	printResults(foundURLs302)
	printResults(foundURLsOther)
}

func printProgress(percentCompleted int) {
	fmt.Printf("\rProgress: %d%%", percentCompleted)
}

func printResults(urls [][]string) {
	if len(urls) == 0 {
		color.Yellow("Bulunan URL yok.")
		return
	}

	fmt.Println("\rProgress: 100%")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Found_URLs", "Status"})

	for _, urlInfo := range urls {
		table.Append(urlInfo)
	}

	table.Render()
}
