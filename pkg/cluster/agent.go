package cluster

import (
	"log"

	"github.com/0x1d/rcond/pkg/config"
	"github.com/hashicorp/serf/serf"
)

type Agent struct {
	Serf *serf.Serf
}

func NewAgent(clusterConfig *config.ClusterConfig) (*Agent, error) {
	config := serf.DefaultConfig()
	config.NodeName = clusterConfig.NodeName
	config.ProtocolVersion = serf.ProtocolVersionMax
	config.MemberlistConfig.SecretKey = []byte(clusterConfig.SecretKey)
	config.MemberlistConfig.AdvertiseAddr = clusterConfig.AdvertiseAddr
	config.MemberlistConfig.AdvertisePort = clusterConfig.AdvertisePort
	config.MemberlistConfig.BindAddr = clusterConfig.BindAddr
	config.MemberlistConfig.BindPort = clusterConfig.BindPort

	serf, err := serf.Create(config)
	if err != nil {
		return nil, err
	}

	return &Agent{Serf: serf}, nil
}

func (a *Agent) Members() ([]serf.Member, error) {
	log.Printf("Getting members of the cluster")
	return a.Serf.Members(), nil
}

func (a *Agent) Join(addrs []string, ignoreOld bool) (int, error) {
	log.Printf("Joining nodes in the cluster: %v", addrs)
	n, err := a.Serf.Join(addrs, ignoreOld)
	if err != nil {
		log.Printf("Failed to join nodes in the cluster: %v", err)
		return 0, err
	}
	log.Printf("Joined %d nodes in the cluster", n)
	return n, nil
}

func (a *Agent) Leave() error {
	return a.Serf.Leave()
}

func (a *Agent) Shutdown() error {
	log.Printf("Shutting down cluster agent")
	return a.Serf.Shutdown()
}
