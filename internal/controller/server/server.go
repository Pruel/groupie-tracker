package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"groupie-tracker/internal/controller/router"
	"groupie-tracker/pkg/config"
)

// Server
type Server struct {
	//http
	httpServer *http.Server
	// ws
	// graphql
}

// New ...
func New(cfg *config.Config, handler *router.Router) *Server {
	slog.Debug("the http server creating")

	defer func() {
		slog.Debug("the http server successfully created")
	}()

	addr := strings.Join([]string{cfg.HTTPServer.Host, cfg.HTTPServer.Port}, ":")
	return &Server{
		httpServer: &http.Server{
			Addr:           addr,
			MaxHeaderBytes: cfg.HTTPServer.MaxHeaderMb << 20, //3 mb
			IdleTimeout:    cfg.HTTPServer.IdleTimeout * time.Second,
			ReadTimeout:    cfg.HTTPServer.ReadTimeout * time.Second,
			WriteTimeout:   cfg.HTTPServer.WriteTimeout * time.Second,
			Handler:        handler.Mux,
		},
	}
}

// Run ...
func (s *Server) Run() error {
	slog.Debug("the http server successful running\n")
	return s.httpServer.ListenAndServe()
}

// Shutdown ...
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
