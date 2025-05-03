// Usage: rcond <address>

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
	if len(os.Args) < 2 {
		usage()
	}

	addr := os.Args[1]
	srv := http.NewServer(addr)
	srv.RegisterRoutes()

	log.Printf("Starting server on %s", addr)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
