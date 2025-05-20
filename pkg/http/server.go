package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/0x1d/rcond/pkg/cluster"
	"github.com/0x1d/rcond/pkg/config"
	"github.com/gorilla/mux"
)

type Server struct {
	router       *mux.Router
	srv          *http.Server
	apiToken     string
	clusterAgent *cluster.Agent
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

func (s *Server) WithClusterAgent(agent *cluster.Agent) *Server {
	s.clusterAgent = agent
	return s
}

func Up(appConfig *config.Config, clusterAgent *cluster.Agent) *Server {
	srv := NewServer(appConfig)
	srv.WithClusterAgent(clusterAgent)
	srv.RegisterRoutes()

	log.Printf("[INFO] Starting API server on %s", appConfig.Rcond.Addr)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
	return srv
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
	s.router.HandleFunc("/network/sta", s.verifyToken(HandleConfigureSTA)).Methods(http.MethodPost)
	s.router.HandleFunc("/network/interface/{interface}", s.verifyToken(HandleNetworkUp)).Methods(http.MethodPut)
	s.router.HandleFunc("/network/interface/{interface}", s.verifyToken(HandleNetworkDown)).Methods(http.MethodDelete)
	s.router.HandleFunc("/network/connection/{uuid}", s.verifyToken(HandleNetworkRemove)).Methods(http.MethodDelete)
	s.router.HandleFunc("/hostname", s.verifyToken(HandleGetHostname)).Methods(http.MethodGet)
	s.router.HandleFunc("/hostname", s.verifyToken(HandleSetHostname)).Methods(http.MethodPost)
	s.router.HandleFunc("/users/{user}/keys", s.verifyToken(HandleAddAuthorizedKey)).Methods(http.MethodPost)
	s.router.HandleFunc("/users/{user}/keys/{fingerprint}", s.verifyToken(HandleRemoveAuthorizedKey)).Methods(http.MethodDelete)
	s.router.HandleFunc("/system/file", s.verifyToken(HandleFileUpload)).Methods(http.MethodPost)
	s.router.HandleFunc("/system/restart", s.verifyToken(HandleReboot)).Methods(http.MethodPost)
	s.router.HandleFunc("/system/shutdown", s.verifyToken(HandleShutdown)).Methods(http.MethodPost)
	s.router.HandleFunc("/cluster/members", s.verifyToken(ClusterAgentHandler(s.clusterAgent, HandleClusterMembers))).Methods(http.MethodGet)
	s.router.HandleFunc("/cluster/join", s.verifyToken(ClusterAgentHandler(s.clusterAgent, HandleClusterJoin))).Methods(http.MethodPost)
	s.router.HandleFunc("/cluster/leave", s.verifyToken(ClusterAgentHandler(s.clusterAgent, HandleClusterLeave))).Methods(http.MethodPost)
	s.router.HandleFunc("/cluster/event", s.verifyToken(ClusterAgentHandler(s.clusterAgent, HandleClusterEvent))).Methods(http.MethodPost)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
