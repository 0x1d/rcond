package http

import (
	"encoding/json"
	"log"
	"net/http"

	network "github.com/0x1d/rcond/pkg/network"
	"github.com/gorilla/mux"
)

type configureAPRequest struct {
	Interface   string `json:"interface"`
	SSID        string `json:"ssid"`
	Password    string `json:"password"`
	Autoconnect bool   `json:"autoconnect"`
}

type configureSTARequest struct {
	Interface   string `json:"interface"`
	SSID        string `json:"ssid"`
	Password    string `json:"password"`
	Autoconnect bool   `json:"autoconnect"`
}

type networkUpRequest struct {
	UUID string `json:"uuid"`
}

type setHostnameRequest struct {
	Hostname string `json:"hostname"`
}

func HandleConfigureSTA(w http.ResponseWriter, r *http.Request) {
	var req configureSTARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Configuring station on interface %s", req.Interface)
	uuid, err := network.ConfigureSTA(req.Interface, req.SSID, req.Password, req.Autoconnect)
	if err != nil {
		log.Printf("Failed to configure station on interface %s: %v", req.Interface, err)
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"uuid": uuid})
}

func HandleConfigureAP(w http.ResponseWriter, r *http.Request) {
	var req configureAPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Configuring access point on interface %s", req.Interface)
	uuid, err := network.ConfigureAP(req.Interface, req.SSID, req.Password, req.Autoconnect)
	if err != nil {
		log.Printf("Failed to configure access point on interface %s: %v", req.Interface, err)
		WriteError(w, err.Error(), http.StatusInternalServerError)
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
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleNetworkUp(w http.ResponseWriter, r *http.Request) {
	var req networkUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	iface := vars["interface"]

	log.Printf("Bringing up network interface %s with UUID %s", iface, req.UUID)
	if err := network.Up(iface, req.UUID); err != nil {
		log.Printf("Failed to bring up network interface %s: %v", iface, err)
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully brought up network interface %s", iface)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleNetworkDown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iface := vars["interface"]
	if err := network.Down(iface); err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleNetworkRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if err := network.Remove(uuid); err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleGetHostname(w http.ResponseWriter, r *http.Request) {
	hostname, err := network.GetHostname()
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"hostname": hostname})
}

func HandleSetHostname(w http.ResponseWriter, r *http.Request) {
	var req setHostnameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := network.SetHostname(req.Hostname); err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
