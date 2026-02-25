package app

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/fx"
)

func NewApp(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("HTTP server starting on %s", server.Addr)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("shutting down server")
			return server.Shutdown(ctx)
		},
	})
}
