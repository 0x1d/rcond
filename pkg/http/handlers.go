package http

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/0x1d/rcond/pkg/network"
	"github.com/0x1d/rcond/pkg/user"
	"github.com/gorilla/mux"
)

const (
	NETWORK_CONNECTION_UUID = "7d706027-727c-4d4c-a816-f0e1b99db8ab"
)

type configureAPRequest struct {
	Interface string `json:"interface"`
	SSID      string `json:"ssid"`
	Password  string `json:"password"`
}

type networkUpRequest struct {
	UUID string `json:"uuid"`
}

type setHostnameRequest struct {
	Hostname string `json:"hostname"`
}

type authorizedKeyRequest struct {
	User   string `json:"user"`
	PubKey string `json:"pubkey"`
}

func HandleConfigureAP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req configureAPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Configuring access point on interface %s", req.Interface)
	uuid, err := network.ConfigureAP(req.Interface, req.SSID, req.Password)
	if err != nil {
		log.Printf("Failed to configure access point on interface %s: %v", req.Interface, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully configured access point on interface %s with UUID %s", req.Interface, uuid)

	resp := struct {
		UUID string `json:"uuid"`
	}{
		UUID: uuid,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func HandleNetworkUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req networkUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	iface := vars["interface"]

	log.Printf("Bringing up network interface %s with UUID %s", iface, req.UUID)
	if err := network.Up(iface, req.UUID); err != nil {
		log.Printf("Failed to bring up network interface %s: %v", iface, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully brought up network interface %s", iface)

	w.WriteHeader(http.StatusOK)
}

func HandleNetworkDown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	iface := vars["interface"]
	if err := network.Down(iface); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleNetworkRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if err := network.Remove(uuid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleGetHostname(w http.ResponseWriter, r *http.Request) {
	hostname, err := network.GetHostname()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(hostname))
}

func HandleSetHostname(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req setHostnameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := network.SetHostname(req.Hostname); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleAddAuthorizedKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authorizedKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	username := vars["user"]

	fingerprint, err := user.AddAuthorizedKey(username, req.PubKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fingerprint))
}

func HandleRemoveAuthorizedKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fingerprint := vars["fingerprint"]
	if fingerprint == "" {
		http.Error(w, "fingerprint parameter is required", http.StatusBadRequest)
		return
	}

	username := vars["user"]
	if username == "" {
		http.Error(w, "user parameter is required", http.StatusBadRequest)
		return
	}

	fingerprintBytes, err := base64.RawURLEncoding.DecodeString(fingerprint)
	if err != nil {
		http.Error(w, "invalid fingerprint base64", http.StatusBadRequest)
		return
	}

	if err := user.RemoveAuthorizedKey(username, string(fingerprintBytes)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
