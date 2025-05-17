// Usage: rcond <address> <api-token>

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/0x1d/rcond/pkg/cluster"
	"github.com/0x1d/rcond/pkg/config"
	http "github.com/0x1d/rcond/pkg/http"
	"github.com/0x1d/rcond/pkg/network"
	"github.com/0x1d/rcond/pkg/system"
	"github.com/godbus/dbus/v5"
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

	configureSystem(appConfig)
	clusterAgent := startClusterAgent(appConfig)
	startApiServer(appConfig, clusterAgent)

	select {}
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

func startApiServer(appConfig *config.Config, clusterAgent *cluster.Agent) *http.Server {
	srv := http.NewServer(appConfig)
	srv.WithClusterAgent(clusterAgent)
	srv.RegisterRoutes()

	log.Printf("[INFO] Starting API server on %s", appConfig.Rcond.Addr)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
	return srv
}

func startClusterAgent(appConfig *config.Config) *cluster.Agent {
	clusterConfig := &appConfig.Cluster
	if clusterConfig.Enabled {
		log.Printf("[INFO] Starting cluster agent on %s:%d", clusterConfig.BindAddr, clusterConfig.BindPort)
		clusterAgent, err := cluster.NewAgent(clusterConfig, cluster.ClusterEventsMap())
		if err != nil {
			log.Fatal(err)
		}
		// join nodes in the cluster if the join addresses are provided
		if len(clusterConfig.Join) > 0 {
			clusterAgent.Join(clusterConfig.Join, true)
		}
		return clusterAgent
	}
	return nil
}

func configureSystem(appConfig *config.Config) error {
	log.Print("[INFO] Configure system")
	// configure hostname
	if err := network.SetHostname(appConfig.Hostname); err != nil {
		log.Printf("[ERROR] %s", err)
	}
	// configure network connections
	for _, connection := range appConfig.Network.Connections {
		err := system.WithDbus(func(conn *dbus.Conn) error {
			_, err := network.AddConnectionWithConfig(conn, &network.ConnectionConfig{
				Type:        connection.Type,
				UUID:        connection.UUID,
				ID:          connection.ID,
				AutoConnect: connection.AutoConnect,
				SSID:        connection.SSID,
				Mode:        connection.Mode,
				Band:        connection.Band,
				Channel:     connection.Channel,
				KeyMgmt:     connection.KeyMgmt,
				PSK:         connection.PSK,
				IPv4Method:  connection.IPv4Method,
				IPv6Method:  connection.IPv6Method,
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}

	}
	log.Print("[INFO] System configured")
	return nil
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
