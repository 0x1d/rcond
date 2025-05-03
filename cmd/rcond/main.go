// Usage: rcond <address> <api-token>

package main

import (
	"fmt"
	"log"
	"os"

	http "github.com/0x1d/rcond/pkg/http"
)

const (
	NETWORK_CONNECTION_UUID = "7d706027-727c-4d4c-a816-f0e1b99db8ab"
)

func usage() {
	fmt.Printf("Usage: %s <address>\n", os.Args[0])
	os.Exit(0)
}

func main() {

	addr := "0.0.0.0:8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	apiToken := os.Getenv("RCOND_API_TOKEN")
	if apiToken == "" {
		log.Fatal("RCOND_API_TOKEN environment variable not set")
	}

	srv := http.NewServer(addr, apiToken)
	srv.RegisterRoutes()

	log.Printf("Starting server on %s", addr)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
