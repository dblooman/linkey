package command

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
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
	filename, _ := filepath.Abs(c.Args().First())
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		color.Red("Unable to load config")
		os.Exit(1)
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
		client := &http.Client{}
		req, err := http.NewRequest("GET", path, nil)

		if err != nil {
			panic(err)
		}
		req.Header.Set("Connection", "close")
		req.Close = true
		go func() {

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()

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
