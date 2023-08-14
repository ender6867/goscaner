


# GoScaner


## Introduction

GoScaner is a command-line tool designed for various scanning tasks. It allows you to perform URL scraping and IP scanning tasks efficiently. Whether you're scanning a website or performing IP-based scans, GoScaner can assist you in identifying relevant information.


## Installation

To use GoScaner, follow these steps:

1. Clone the repository: `git clone https://github.com/ender6867/goscaner.git`
2. Navigate to the project directory: `cd goscaner`
3. Build the binary: `go build`


## Usage

GoScaner supports the following flags:

- `-h, --help`: Displays the help message to guide you through usage.
- `-s, --scan-ip`: Perform IP scan.
- `-u, --url`: Specifies the URL to be scraped. Example: `-u https://www.example.com`
- `-w, --wordlist`: Path to the wordlist file. Example: `-w path/wordlist.txt`

## Configuration

Before running GoScaner, make sure to set up your API key:

1. Create a `.env` file in the project directory.
2. Add your API key to the `.env` file:

## Examples

Here are some example usages to perform scanning and see the results:

1. Perform a URL scan:
   ```sh
   ./goscaner -u https://www.example.com -w path/wordlist.txt
2. Perform an IP scan:
      ```sh
   ./goscaner -s -u https://www.example.com -w path/wordlist.txt

