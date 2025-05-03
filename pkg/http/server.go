package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router   *mux.Router
	srv      *http.Server
	apiToken string
}

func NewServer(addr string, apiToken string) *Server {
	router := mux.NewRouter()

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &Server{
		router:   router,
		srv:      srv,
		apiToken: apiToken,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) verifyToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-API-Token")
		if token == "" || token != s.apiToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/health", s.healthHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/network/up", s.verifyToken(HandleNetworkUp)).Methods(http.MethodPost)
	s.router.HandleFunc("/network/down", s.verifyToken(HandleNetworkDown)).Methods(http.MethodPost)
	s.router.HandleFunc("/network/remove", s.verifyToken(HandleNetworkRemove)).Methods(http.MethodPost)
	s.router.HandleFunc("/hostname", s.verifyToken(HandleGetHostname)).Methods(http.MethodGet)
	s.router.HandleFunc("/hostname", s.verifyToken(HandleSetHostname)).Methods(http.MethodPost)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
