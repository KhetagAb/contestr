// cmd/api/main.go
package main

import (
	"contestr/internal/configs"
	httphandlers "contestr/internal/handlers/http"
	"contestr/internal/services/codeforces"
	"contestr/internal/transport"
	"contestr/pkg/config"
	"contestr/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, ctx := config.InitConfigAndContext()
	httpServer := createHttpServer(ctx, cfg)
	startServer(ctx, httpServer)
	awaitGracefulShutdown(ctx, cfg, httpServer)
}

func startServer(ctx context.Context, httpServer *transport.Server) {
	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logger.Errorf(ctx, "failed to start http server: %v", err)
		}
	}()
}

func createHttpServer(ctx context.Context, cfg *configs.Config) *transport.Server {
	cfService := codeforces.NewService(cfg)
	handlers := httphandlers.NewHandlers(cfService)

	httpServer := transport.NewHTTPServer(ctx, cfg)
	httpServer.RegisterHandlers(handlers)
	return httpServer
}

func awaitGracefulShutdown(ctx context.Context, cfg *configs.Config, httpServer *transport.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		logger.Fatalf(ctx, "server forced to shutdown: %v", err)
	}

	logger.Info(ctx, "server exited properly")
}
