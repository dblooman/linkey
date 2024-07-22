package command

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type Config struct {
	Base        string
	Concurrency int
	Headers     []struct {
		Key   string
		Value string
	}
	Paths      []string
	Statuscode int
}

func CmdCheck(c *cli.Context) {
	filename, err := filepath.Abs(c.Args().First())
	if err != nil {
		color.Red("Failed to get absolute path: %v", err)
		os.Exit(1)
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		color.Red("Unable to load config: %v", err)
		os.Exit(1)
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		color.Red("Failed to unmarshal YAML: %v", err)
		os.Exit(1)
	}

	if err := validateConfig(config); err != nil {
		color.Red("Invalid config: %v", err)
		os.Exit(1)
	}

	response := checker(config)

	totalUrls := color.New(color.FgWhite).PrintfFunc()
	totalUrls("\nTotal URLs checked: %d", len(config.Paths))

	if len(response) > 0 {
		showBroken(response)
	} else {
		color.Green("\n\nAll URLs are good buddy :)")
	}
}

func validateConfig(config Config) error {
	if config.Base == "" {
		return fmt.Errorf("base URL cannot be empty")
	}
	if config.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0")
	}
	if config.Statuscode == 0 {
		return fmt.Errorf("status code must be specified")
	}
	if len(config.Paths) == 0 {
		return fmt.Errorf("no paths specified")
	}
	return nil
}

func checker(config Config) []string {
	var brokenUrls []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	headers := make(map[string]string)
	for _, element := range config.Headers {
		headers[element.Key] = element.Value
	}

	sem := make(chan struct{}, config.Concurrency)
	client := &http.Client{}

	for _, path := range config.Paths {
		wg.Add(1)

		go func(path string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			if checkURL(client, config.Base+path, headers, config.Statuscode) {
				mu.Lock()
				brokenUrls = append(brokenUrls, config.Base+path)
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()
	return brokenUrls
}

func checkURL(client *http.Client, url string, headers map[string]string, expectedStatusCode int) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red("Failed to create request for %s: %v", url, err)
		return false
	}
	req.Header.Set("Connection", "close")
	req.Close = true

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		color.Red("Failed to make request to %s: %v", url, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == expectedStatusCode {
		color.Green("Status is good for: %s", url)
		return false
	} else {
		printBrokenUrls(url, resp.StatusCode)
		return true
	}
}

func printBrokenUrls(path string, code int) {
	red := color.New(color.FgRed).PrintfFunc()
	red("Error with: %s, status is %d\n", path, code)
}

func showBroken(response []string) {
	color.Red("\n\nYou got a broken link Buddy :(")

	for _, element := range response {
		color.Yellow(element)
	}

	os.Exit(1)
}
