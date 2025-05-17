package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/0x1d/rcond/pkg/user"
	"github.com/gorilla/mux"
)

type authorizedKeyRequest struct {
	User   string `json:"user"`
	PubKey string `json:"pubkey"`
}

func HandleAddAuthorizedKey(w http.ResponseWriter, r *http.Request) {
	var req authorizedKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	username := vars["user"]

	fingerprint, err := user.AddAuthorizedKey(username, req.PubKey)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"fingerprint": fingerprint})
}

func HandleRemoveAuthorizedKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fingerprint := vars["fingerprint"]
	if fingerprint == "" {
		WriteError(w, "fingerprint parameter is required", http.StatusBadRequest)
		return
	}

	username := vars["user"]
	if username == "" {
		WriteError(w, "user parameter is required", http.StatusBadRequest)
		return
	}

	fingerprintBytes, err := base64.RawURLEncoding.DecodeString(fingerprint)
	if err != nil {
		WriteError(w, "invalid fingerprint base64", http.StatusBadRequest)
		return
	}

	if err := user.RemoveAuthorizedKey(username, string(fingerprintBytes)); err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
