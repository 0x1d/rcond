package rcond

import (
	"log"

	"github.com/0x1d/rcond/pkg/cluster"
	"github.com/0x1d/rcond/pkg/config"
	"github.com/0x1d/rcond/pkg/http"
	"github.com/0x1d/rcond/pkg/system"
)

type Node struct {
	Config       *config.Config
	ClusterAgent *cluster.Agent
	HttpApi      *http.Server
}

func NewNode(appConfig *config.Config) *Node {
	return &Node{
		Config:       appConfig,
		HttpApi:      Api(appConfig),
		ClusterAgent: Cluster(&appConfig.Cluster),
	}
}

func (n *Node) Up() {
	system.Configure(n.Config)
	n.HttpApi.WithClusterAgent(n.ClusterAgent)
	n.HttpApi.RegisterRoutes()

	log.Printf("[INFO] Starting API server on %s", n.Config.Rcond.Addr)
	if err := n.HttpApi.Start(); err != nil {
		log.Fatal(err)
	}
}

func Api(appConfig *config.Config) *http.Server {
	srv := http.NewServer(appConfig)
	return srv
}

func Cluster(clusterConfig *config.ClusterConfig) *cluster.Agent {
	if clusterConfig.Enabled {
		log.Printf("[INFO] Starting cluster agent on %s:%d", clusterConfig.BindAddr, clusterConfig.BindPort)
		clusterAgent, err := cluster.NewAgent(clusterConfig, cluster.ClusterEventsMap())
		if err != nil {
			log.Print(err)
			return nil
		}
		// join nodes in the cluster if the join addresses are provided
		if len(clusterConfig.Join) > 0 {
			clusterAgent.Join(clusterConfig.Join, true)
		}
		return clusterAgent
	}
	return nil
}
