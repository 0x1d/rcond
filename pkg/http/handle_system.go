package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/0x1d/rcond/pkg/system"
)

func HandleReboot(w http.ResponseWriter, r *http.Request) {
	if err := system.Restart(); err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleShutdown(w http.ResponseWriter, r *http.Request) {
	if err := system.Shutdown(); err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var fileUpload struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&fileUpload); err != nil {
		WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Decode base64 encoded content to bytes
	contentBytes, err := base64.StdEncoding.DecodeString(fileUpload.Content)
	if err != nil {
		WriteError(w, "Failed to decode base64 content", http.StatusBadRequest)
		return
	}

	// Store the file
	if err := system.StoreFile(fileUpload.Path, contentBytes); err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
