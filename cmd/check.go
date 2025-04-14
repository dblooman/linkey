package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/dblooman/linkey/internal/linkey"
	"github.com/dblooman/linkey/internal/models"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		CmdCheck(args)
	},
}

var filePath string

func init() {
	checkCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the configuration file")
	rootCmd.AddCommand(checkCmd)
}

func CmdCheck(args []string) {
	var filename string

	if filePath != "" {
		filename = filePath
	} else if len(args) > 0 {
		filename = args[0]
	} else {
		logErrorAndExit(fmt.Errorf("no configuration file provided"))
	}

	config := loadConfig(filename)

	request := linkey.New(config)
	response := request.Checker(config)

	totalUrls := color.New(color.FgWhite).PrintfFunc()
	totalUrls("\nTotal URLs checked: %d", len(config.Paths))

	if len(response) > 0 {
		showBroken(response)
	} else {
		color.Green("\n\nAll URLs are good buddy :)")
	}
}

func validateConfig(config models.Config) error {
	if config.Base == "" {
		return fmt.Errorf("base URL cannot be empty")
	}
	if config.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0")
	}
	if config.StatusCode == 0 {
		return fmt.Errorf("status code must be specified")
	}
	if len(config.Paths) == 0 {
		return fmt.Errorf("no paths specified")
	}
	return nil
}

func logErrorAndExit(err error) {
	color.Red("Error: %v", err)
	os.Exit(1)
}

func showBroken(response []string) {
	color.Red("\n\nYou got a broken link Buddy :(")

	for _, element := range response {
		color.Yellow(element)
	}
}

func loadConfig(filename string) models.Config {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		logErrorAndExit(fmt.Errorf("failed to get absolute path: %w", err))
	}

	yamlFile, err := os.ReadFile(absPath)
	if err != nil {
		logErrorAndExit(fmt.Errorf("unable to load config: %w", err))
	}

	var config models.Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		logErrorAndExit(fmt.Errorf("failed to unmarshal YAML: %w", err))
	}

	if err := validateConfig(config); err != nil {
		logErrorAndExit(fmt.Errorf("invalid config: %w", err))
	}

	return config
}
