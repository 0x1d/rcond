// Usage: rcond <address> <api-token>

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/0x1d/rcond/pkg/config"
	http "github.com/0x1d/rcond/pkg/http"
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

	srv := http.NewServer(appConfig)
	srv.RegisterRoutes()

	log.Printf("Starting server on %s", appConfig.Rcond.Addr)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() (*config.Config, error) {
	configPath := "/etc/rcond/config.yaml"
	appConfig := &config.Config{}
	help := false

	flag.StringVar(&configPath, "config", configPath, "Path to the configuration file")
	flag.StringVar(&appConfig.Rcond.Addr, "addr", "", "Address to bind the HTTP server to")
	flag.StringVar(&appConfig.Rcond.ApiToken, "token", "", "API token to use for authentication")
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

	// Override config values from environment variables and flags
	overrideConfigValuesFromEnv(map[string]*string{
		"RCOND_ADDR":      &appConfig.Rcond.Addr,
		"RCOND_API_TOKEN": &appConfig.Rcond.ApiToken,
	})

	overrideConfigValuesFromFlag(map[string]*string{
		"addr":  &appConfig.Rcond.Addr,
		"token": &appConfig.Rcond.ApiToken,
	})

	// Validate required fields
	if err := validateRequiredFields(map[string]*string{
		"addr":  &appConfig.Rcond.Addr,
		"token": &appConfig.Rcond.ApiToken,
	}); err != nil {
		return nil, err
	}

	return appConfig, nil
}

func overrideConfigValuesFromEnv(envMap map[string]*string) {
	for varName, configValue := range envMap {
		if envValue, ok := os.LookupEnv(varName); ok {
			*configValue = envValue
		}
	}
}

func overrideConfigValuesFromFlag(flagMap map[string]*string) {
	for flagName, configValue := range flagMap {
		if flagValue := flag.Lookup(flagName).Value.String(); flagValue != "" {
			*configValue = flagValue
		}
	}
}

func validateRequiredFields(fields map[string]*string) error {
	for name, value := range fields {
		if *value == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	return nil
}
