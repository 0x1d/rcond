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

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func HandleConfigureAP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req configureAPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Configuring access point on interface %s", req.Interface)
	uuid, err := network.ConfigureAP(req.Interface, req.SSID, req.Password)
	if err != nil {
		log.Printf("Failed to configure access point on interface %s: %v", req.Interface, err)
		writeError(w, err.Error(), http.StatusInternalServerError)
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
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleNetworkUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req networkUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	iface := vars["interface"]

	log.Printf("Bringing up network interface %s with UUID %s", iface, req.UUID)
	if err := network.Up(iface, req.UUID); err != nil {
		log.Printf("Failed to bring up network interface %s: %v", iface, err)
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully brought up network interface %s", iface)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleNetworkDown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	iface := vars["interface"]
	if err := network.Down(iface); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleNetworkRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if err := network.Remove(uuid); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleGetHostname(w http.ResponseWriter, r *http.Request) {
	hostname, err := network.GetHostname()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"hostname": hostname})
}

func HandleSetHostname(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req setHostnameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := network.SetHostname(req.Hostname); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleAddAuthorizedKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authorizedKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	username := vars["user"]

	fingerprint, err := user.AddAuthorizedKey(username, req.PubKey)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"fingerprint": fingerprint})
}

func HandleRemoveAuthorizedKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fingerprint := vars["fingerprint"]
	if fingerprint == "" {
		writeError(w, "fingerprint parameter is required", http.StatusBadRequest)
		return
	}

	username := vars["user"]
	if username == "" {
		writeError(w, "user parameter is required", http.StatusBadRequest)
		return
	}

	fingerprintBytes, err := base64.RawURLEncoding.DecodeString(fingerprint)
	if err != nil {
		writeError(w, "invalid fingerprint base64", http.StatusBadRequest)
		return
	}

	if err := user.RemoveAuthorizedKey(username, string(fingerprintBytes)); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
