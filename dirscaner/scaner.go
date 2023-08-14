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

	var foundURLs []string
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

			if resp.StatusCode == 200 {
				mu.Lock()
				foundURLs = append(foundURLs, targetURL)
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

	printResults(foundURLs)
}

func printProgress(percentCompleted int) {
	fmt.Printf("\rProgress: %d%%", percentCompleted)
}

func printResults(urls []string) {
	if len(urls) == 0 {
		color.Yellow("Bulunan URL yok.")
		return
	}

	fmt.Println("\rProgress: 100%")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Found_URLs"})

	for _, url := range urls {
		table.Append([]string{color.RedString(url)})
	}

	table.Render()
}
