package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
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

func main() {

	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	response := checker(config)

	totalUrls := color.New(color.FgWhite).PrintfFunc()
	totalUrls("\nTotal URL's checked: %d", len(config.Paths))

	if len(response) > 0 {
		showBroken(response)
	} else {
		color.Green("\n\nAll URL's are good buddy :)")
	}

}

func checker(config Config) []string {

	var brokenUrls []string

	var wg sync.WaitGroup

	for _, element := range config.Paths {
		path := config.Base + element
		wg.Add(1)
		go func() {
			resp, err := http.Get(path)

			if err != nil {
				panic(err)
			}

			if resp.StatusCode == config.Statuscode {
				green := color.New(color.FgGreen).PrintfFunc()
				green("\nStatus is good for: %s", path)
			} else {
				brokenUrls = append(brokenUrls, path)
				printBrokenUrls(path, resp.StatusCode)
			}

			wg.Done()
		}()
	}
	wg.Wait()

	return brokenUrls
}

func printBrokenUrls(path string, code int) {
	red := color.New(color.FgRed).PrintfFunc()
	red("\nError with: %s", path)
	red("\nStatus is: %d", code)

}

func showBroken(response []string) {
	color.Red("\n\nYou got a broken link Buddy :(")

	for _, element := range response {
		color.Yellow(element)
	}

	os.Exit(1)
}
