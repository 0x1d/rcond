package cluster

import (
	"log"

	"github.com/0x1d/rcond/pkg/network"
	"github.com/0x1d/rcond/pkg/system"
)

func ClusterEventsMap() map[string]func([]byte) {
	return map[string]func([]byte){
		"printHostname": printHostname,
		"restart":       restart,
		"shutdown":      shutdown,
	}
}

func restart(payload []byte) {
	if err := system.Restart(); err != nil {
		log.Printf("[ERROR] (ClusterEvent:restart) failed: %s", err)
	}
}

func shutdown(payload []byte) {
	if err := system.Shutdown(); err != nil {
		log.Printf("[ERROR] (ClusterEvent:shutdown) failed: %s", err)
	}
}

// just a sample function to test event functionality
func printHostname(payload []byte) {
	hostname, _ := network.GetHostname()
	log.Printf("[INFO] (ClusterEvent:printHostname): %s", hostname)
}
