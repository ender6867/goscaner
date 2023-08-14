package main

import (
	"bufio"
	dirscanner "goscaner/dirscaner"
	"goscaner/lookup"
	"goscaner/scraper"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var url string
	var wordlistPath string
	var scanIP bool

	var rootCmd = &cobra.Command{
		Use:   "goscaner",
		Short: "wordpress scanner",
		Long: `                                            
  __ _  ___  ___  ___ __ _ _ __   ___ _ __ 
 / _` + "`" + ` |/ _ \/ __|/ __/ _` + "`" + ` | '_ \ / _ \ '__|
| (_| | (_) \__ \ (_| (_| | | | |  __/ |   
 \__, |\___/|___/\___\__,_|_| |_|\___|_|   
 |___/`,
		Example: "goscaner -u https://www.example.com -w path/wordlist.txt -s",

		Run: func(cmd *cobra.Command, args []string) {
			var wordlistContent []string
			if wordlistPath != "" {
				var err error
				wordlistContent, err = readWordlist(wordlistPath)
				if err != nil {
					log.Fatal(err)
				}
			}
			scraper.Scrape(url)
			if wordlistPath != "" {
				dirscanner.Scan(url, wordlistContent)
			}
			if scanIP {
				lookup.Scan(url)
			}
		},
	}

	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL to scrape")
	rootCmd.Flags().StringVarP(&wordlistPath, "wordlist", "w", "", "Path to wordlist file")
	rootCmd.Flags().BoolVarP(&scanIP, "scan-ip", "s", false, "Perform IP scan")
	rootCmd.MarkFlagRequired("url")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func readWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return content, nil
}
