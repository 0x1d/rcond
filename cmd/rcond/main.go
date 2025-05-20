// Usage: rcond <address> <api-token>

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0x1d/rcond/pkg/config"
	"github.com/0x1d/rcond/pkg/rcond"
)

func usage() {
	fmt.Println("Usage: rcond <flags>")
	flag.PrintDefaults()
}

func main() {
	appConfig, err := loadConfig()
	if err != nil {
		usage()
		fmt.Printf("\nFailed to load config: %v\n", err)
		os.Exit(1)
	}

	rcond.NewNode(appConfig).Up()

	select {}
}

func loadConfig() (*config.Config, error) {
	configPath := "/etc/rcond/config.yaml"
	appConfig := &config.Config{}
	help := false

	flag.StringVar(&configPath, "config", configPath, "Path to the configuration file")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	// Load config from file
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		configFile, err := config.LoadConfig(configPath)
		if err != nil {
			return nil, err
		}
		appConfig = configFile
	}

	// Validate required fields
	if err := validateRequiredFields(map[string]*string{
		"addr":  &appConfig.Rcond.Addr,
		"token": &appConfig.Rcond.ApiToken,
	}); err != nil {
		return nil, err
	}

	return appConfig, nil
}

func validateRequiredFields(fields map[string]*string) error {
	for name, value := range fields {
		if *value == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	return nil
}
