package http

import (
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
