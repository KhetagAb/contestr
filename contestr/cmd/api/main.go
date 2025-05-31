// cmd/api/main.go
package main

import (
	"contestr/internal/config"
	"contestr/internal/transport"
	"contestr/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig("./configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}

	logger.Init(cfg.App.Name, "info")
	ctx := context.Background()

	httpServer := createHttpServer(ctx, cfg)

	logger.Info(ctx, "starting http server")
	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logger.Errorf(ctx, "failed to start http server: %v", err)
		}
	}()

	awaitGracefulShutdown(ctx, cfg, httpServer)
}

func createHttpServer(ctx context.Context, cfg *config.Config) *transport.Server {
	httpServer := transport.NewHTTPServer(ctx, cfg)
	httpServer.RegisterHandlers()
	return httpServer
}

func awaitGracefulShutdown(ctx context.Context, cfg *config.Config, httpServer *transport.Server) {
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
