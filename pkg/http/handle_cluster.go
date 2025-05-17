package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/0x1d/rcond/pkg/cluster"
)

type clusterEventRequest struct {
	Name    string `json:"name"`
	Payload string `json:"payload,omitempty"`
}

func ClusterAgentHandler(agent *cluster.Agent, handler func(http.ResponseWriter, *http.Request, *cluster.Agent)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, agent)
	}
}

func HandleClusterJoin(w http.ResponseWriter, r *http.Request, agent *cluster.Agent) {
	var joinRequest struct {
		Join []string `json:"join"`
	}
	err := json.NewDecoder(r.Body).Decode(&joinRequest)
	if err != nil {
		WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	joinAddrs := strings.Join(joinRequest.Join, ",")
	if joinAddrs == "" {
		WriteError(w, "No join addresses provided", http.StatusBadRequest)
		return
	}

	addrs := strings.Split(joinAddrs, ",")
	n, err := agent.Join(addrs, true)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"joined": n})
}

func HandleClusterLeave(w http.ResponseWriter, r *http.Request, agent *cluster.Agent) {
	err := agent.Leave()
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleClusterMembers(w http.ResponseWriter, r *http.Request, agent *cluster.Agent) {
	if agent == nil {
		WriteError(w, "cluster agent is not initialized", http.StatusInternalServerError)
		return
	}
	members, err := agent.Members()
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func HandleClusterEvent(w http.ResponseWriter, r *http.Request, agent *cluster.Agent) {
	if agent == nil {
		WriteError(w, "cluster agent is not initialized", http.StatusInternalServerError)
		return
	}
	var req clusterEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	event := cluster.ClusterEvent{
		Name: req.Name,
		Data: []byte(req.Payload),
	}
	err := agent.Event(event)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
