package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	srv    *http.Server
}

func NewServer(addr string) *Server {
	router := mux.NewRouter()

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &Server{
		router: router,
		srv:    srv,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/health", s.healthHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/network/up", HandleNetworkUp).Methods(http.MethodPost)
	s.router.HandleFunc("/network/down", HandleNetworkDown).Methods(http.MethodPost)
	s.router.HandleFunc("/network/remove", HandleNetworkRemove).Methods(http.MethodPost)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
