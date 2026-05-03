package server

import (
	"hardware_store/internal/config"
	"net/http"
)

func NewServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: 0,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
}
