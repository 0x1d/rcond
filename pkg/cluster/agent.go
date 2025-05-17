package cluster

import (
	"encoding/json"
	"log"
	"os"

	"github.com/0x1d/rcond/pkg/config"
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/serf/serf"
)

// Agent represents a Serf cluster agent.
type Agent struct {
	Serf *serf.Serf
}

// ClusterEvent represents a custom event that will be sent to the Serf cluster.
type ClusterEvent struct {
	Name string
	Data []byte
}

// NewAgent creates a new Serf cluster agent with the given configuration and event handlers.
func NewAgent(clusterConfig *config.ClusterConfig, clusterEvents map[string]func([]byte)) (*Agent, error) {
	config := serf.DefaultConfig()
	config.Init()
	logFilter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(clusterConfig.LogLevel),
		Writer:   os.Stderr,
	}
	config.LogOutput = logFilter
	config.MemberlistConfig.LogOutput = logFilter
	config.NodeName = clusterConfig.NodeName
	config.ProtocolVersion = serf.ProtocolVersionMax
	config.MemberlistConfig.SecretKey = []byte(clusterConfig.SecretKey)
	config.MemberlistConfig.AdvertiseAddr = clusterConfig.AdvertiseAddr
	config.MemberlistConfig.AdvertisePort = clusterConfig.AdvertisePort
	config.MemberlistConfig.BindAddr = clusterConfig.BindAddr
	config.MemberlistConfig.BindPort = clusterConfig.BindPort

	// Setup event channel
	eventCh := make(chan serf.Event, 10)
	config.EventCh = eventCh
	go handleEvents(eventCh, clusterEvents)

	// Start Serf
	serf, err := serf.Create(config)
	if err != nil {
		return nil, err
	}

	return &Agent{Serf: serf}, nil
}

// Event sends a custom event to the Serf cluster.
// It marshals the provided ClusterEvent into JSON and then uses Serf's UserEvent method to send the event.
func (a *Agent) Event(event ClusterEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := a.Serf.UserEvent(event.Name, eventData, false); err != nil {
		return err
	}
	return nil
}

// Members returns a list of members in the Serf cluster.
func (a *Agent) Members() ([]serf.Member, error) {
	log.Printf("[INFO] Getting members of the cluster")
	return a.Serf.Members(), nil
}

// Join attempts to join the Serf cluster with the given addresses, optionally ignoring old nodes.
func (a *Agent) Join(addrs []string, ignoreOld bool) (int, error) {
	log.Printf("[INFO] Joining nodes in the cluster: %v", addrs)
	n, err := a.Serf.Join(addrs, ignoreOld)
	if err != nil {
		log.Printf("[ERROR] Failed to join nodes in the cluster: %v", err)
		return 0, err
	}
	log.Printf("[INFO] Joined %d nodes in the cluster", n)
	return n, nil
}

// Leave causes the agent to leave the Serf cluster.
func (a *Agent) Leave() error {
	return a.Serf.Leave()
}

// Shutdown shuts down the Serf cluster agent.
func (a *Agent) Shutdown() error {
	log.Printf("[INFO] Shutting down cluster agent")
	return a.Serf.Shutdown()
}

// handleEvents handles Serf events received on the event channel.
func handleEvents(eventCh chan serf.Event, clusterEvents map[string]func([]byte)) {
	eventHandlers := clusterEvents
	for event := range eventCh {
		switch event.EventType() {
		case serf.EventUser:
			userEvent := event.(serf.UserEvent)
			if handler, ok := eventHandlers[userEvent.Name]; ok {
				handler(userEvent.Payload)
			} else {
				log.Printf("[INFO] No event handler found for event: %s", userEvent.Name)
			}
		default:
			log.Printf("[INFO] Received event: %s\n", event.EventType())
		}
	}
}
