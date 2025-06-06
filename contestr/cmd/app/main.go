package main

import (
	"contestr/internal/configs"
	"contestr/internal/generated/app"
	"contestr/internal/transport"
	"contestr/pkg/logger"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	svc, err := app.InitializeService()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize server: %v", err))
	}
	cfg := svc.Cfg
	ctx := svc.Ctx

	server := svc.HttpServer
	go func() {
		if err := server.Start(ctx); err != nil {
			logger.Errorf(ctx, "failed to start http server: %v", err)
		}
	}()

	bot := svc.TgBot
	go func() {
		if err := bot.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf(ctx, "failed to start bot: %v", err)
		}
	}()

	awaitGracefulShutdown(ctx, cfg, bot, server)
}

func awaitGracefulShutdown(ctx context.Context, cfg *configs.Config, bot *transport.TgBot, httpServer *transport.HTTPServer) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "shutting down bot...")
	bot.Stop(ctx)

	logger.Info(ctx, "shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		logger.Fatalf(ctx, "server forced to shutdown: %v", err)
	}
	logger.Info(ctx, "server exited properly")
}
