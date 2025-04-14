package linkey

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dblooman/linkey/internal/models"
	"github.com/fatih/color"
)

type Linkey struct {
	client *http.Client
	config models.Config
}

func New(config models.Config) *Linkey {
	return &Linkey{
		client: &http.Client{},
		config: config,
	}
}

func (l *Linkey) checkURL(ctx context.Context, url string, headers map[string]string, expectedStatusCode int) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logError(fmt.Errorf("failed to create request for %s: %w", url, err))
		return false
	}
	req.Header.Set("Connection", "close")
	req.Close = true

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		logError(fmt.Errorf("failed to make request to %s: %w", url, err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == expectedStatusCode {
		color.Green("Status is good for: %s", url)
		return true
	} else {
		printBrokenUrls(url, resp.StatusCode)
		return false
	}
}

func (l *Linkey) Checker(config models.Config) []string {
	var wg sync.WaitGroup
	errorUrls := make(chan string, len(config.Paths))
	sem := make(chan struct{}, config.Concurrency)

	headers := make(map[string]string)
	for _, element := range config.Headers {
		headers[element.Key] = element.Value
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, path := range config.Paths {
		wg.Add(1)

		go func(path string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			if !l.checkURL(ctx, config.Base+path, headers, config.StatusCode) {
				errorUrls <- config.Base + path
			}
		}(path)
	}

	wg.Wait()
	close(errorUrls)

	var brokenUrls []string
	for url := range errorUrls {
		brokenUrls = append(brokenUrls, url)
	}

	return brokenUrls
}

func printBrokenUrls(path string, code int) {
	red := color.New(color.FgRed).PrintfFunc()
	red("Error with: %s, Status code: %d\n", path, code)
}

func logError(err error) {
	color.Red("Error: %v", err)
}
