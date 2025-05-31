package main

import (
	"contestr/cmd/api"
	"contestr/cmd/tgbot"
	"contestr/internal/configs"
	"contestr/internal/transport"
	"contestr/pkg/config"
	"contestr/pkg/logger"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, ctx := config.InitConfigAndContext()

	httpServer := api.CreateHttpServer(ctx, cfg)
	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logger.Errorf(ctx, "failed to start http server: %v", err)
		}
	}()

	bot := tgbot.CreateBot(cfg)
	go func() {
		if err := bot.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf(ctx, "failed to start bot: %v", err)
		}
	}()

	awaitGracefulShutdown(ctx, cfg, bot, httpServer)
}

func awaitGracefulShutdown(ctx context.Context, cfg *configs.Config, bot *transport.Bot, httpServer *transport.Server) {
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
