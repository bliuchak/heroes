package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bliuchak/heroes/internal/config"
	"github.com/bliuchak/heroes/internal/server/middleware"
	"github.com/bliuchak/heroes/internal/storage"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

// Serverer is interface for server
type Serverer interface {
	InitRouter()
	SetMiddleware()
	Run() error
	Stop(timeout time.Duration) error
}

// Server app container for main dependencies
type Server struct {
	Router  *mux.Router
	Storage storage.Storager
	Logger  zerolog.Logger
	Config  config.Config

	httpServer *http.Server
}

// NewServer create pointer for new server structure
func NewServer(storage storage.Storager, logger zerolog.Logger, config config.Config) *Server {
	return &Server{
		Storage: storage,
		Logger:  logger,
		Config:  config,
	}
}

// InitRouter initialize router
func (s *Server) InitRouter() {
	s.Router = mux.NewRouter()
}

// SetMiddleware middleware setter
func (s *Server) SetMiddleware() {
	md := middleware.Middleware{Logger: s.Logger}
	md.SetLogger(s.Logger)

	s.Router.Use(md.HTTPLogger)
}

// Run runs http server
func (s *Server) Run() error {
	s.InitRouter()
	s.SetRoutes()

	s.SetMiddleware()

	s.httpServer = &http.Server{
		Handler:      s.Router,
		Addr:         ":" + strconv.Itoa(s.Config.Server.Port),
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

// Stop shutdown http server with timeout
func (s *Server) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
