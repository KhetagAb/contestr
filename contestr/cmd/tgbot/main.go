package main

import (
	"contestr/internal/config"
	"contestr/internal/delivery"
	"contestr/internal/handlers/tgbot"
	"contestr/internal/usecases"
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
	logger.Init("contestr", "info")

	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.LoadConfig("./configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	bot := createBot(err, cfg)

	go func() {
		if err := bot.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf(ctx, "bot error: %v", err)
		}
		logger.Info(ctx, "telegram bot started")
	}()

	awaitShutdown(ctx, bot, cancel)
}

func createBot(err error, cfg *config.Config) *delivery.Bot {
	botUseCase := usecases.NewBotUseCase()

	botHandlers := tgbot.NewHandlers(botUseCase)
	bot, err := delivery.NewBot(cfg, botHandlers)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize Telegram bot: %v", err))
	}
	return bot
}

func awaitShutdown(ctx context.Context, bot *delivery.Bot, cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "shutting down bot...")

	bot.Stop(ctx)
	cancel()
}
