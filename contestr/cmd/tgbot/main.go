package main

import (
	"contestr/internal/configs"
	"contestr/internal/handlers/tgbot"
	"contestr/internal/transport"
	"contestr/internal/usecases"
	"contestr/pkg/config"
	"contestr/pkg/logger"
	"context"
	"errors"
	"fmt"
	_ "go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, ctx := config.InitConfigAndContext()

	bot := createBot(cfg)
	logger.Info(ctx, "starting telegram bot...")
	go func() {
		if err := bot.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf(ctx, "bot error: %v", err)
		}
	}()
	logger.Info(ctx, "telegram bot started")

	awaitShutdown(ctx, bot)
}

func createBot(cfg *configs.Config) *transport.Bot {
	botUseCase := usecases.NewBotUseCase()

	botHandlers := tgbot.NewHandlers(botUseCase)
	bot, err := transport.NewBot(cfg, botHandlers)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize Telegram bot: %v", err))
	}
	bot.RegisterHandlers()
	return bot
}

func awaitShutdown(ctx context.Context, bot *transport.Bot) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "shutting down bot...")

	bot.Stop(ctx)
}
