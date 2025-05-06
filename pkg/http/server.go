package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/0x1d/rcond/pkg/config"
	"github.com/gorilla/mux"
)

type Server struct {
	router   *mux.Router
	srv      *http.Server
	apiToken string
}

func NewServer(cfg *config.Config) *Server {
	if cfg.Rcond.Addr == "" || cfg.Rcond.ApiToken == "" {
		panic("addr or api_token is not set")
	}

	router := mux.NewRouter()

	srv := &http.Server{
		Addr:         cfg.Rcond.Addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &Server{
		router:   router,
		srv:      srv,
		apiToken: cfg.Rcond.ApiToken,
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
	s.router.HandleFunc("/network/ap", s.verifyToken(HandleConfigureAP)).Methods(http.MethodPost)
	s.router.HandleFunc("/network/interface/{interface}", s.verifyToken(HandleNetworkUp)).Methods(http.MethodPut)
	s.router.HandleFunc("/network/interface/{interface}", s.verifyToken(HandleNetworkDown)).Methods(http.MethodDelete)
	s.router.HandleFunc("/network/connection/{uuid}", s.verifyToken(HandleNetworkRemove)).Methods(http.MethodDelete)
	s.router.HandleFunc("/hostname", s.verifyToken(HandleGetHostname)).Methods(http.MethodGet)
	s.router.HandleFunc("/hostname", s.verifyToken(HandleSetHostname)).Methods(http.MethodPost)
	s.router.HandleFunc("/users/{user}/keys", s.verifyToken(HandleAddAuthorizedKey)).Methods(http.MethodPost)
	s.router.HandleFunc("/users/{user}/keys/{fingerprint}", s.verifyToken(HandleRemoveAuthorizedKey)).Methods(http.MethodDelete)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
