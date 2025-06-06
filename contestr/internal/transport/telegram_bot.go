package transport

import (
	"contestr/internal/configs"
	"contestr/internal/handlers/tgbot"
	"contestr/pkg/logger"
	"context"
	"github.com/go-telegram/bot"
)

type TgBot struct {
	bot      *bot.Bot
	cfg      *configs.Config
	handlers *tgbot.Handlers
}

func NewBot(cfg *configs.Config, handlers *tgbot.Handlers) (*TgBot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandle),
		bot.WithDebug(),
	}

	b, err := bot.New(cfg.Telegram.Token, opts...)
	if err != nil {
		return nil, err
	}

	tgBot := TgBot{
		bot:      b,
		cfg:      cfg,
		handlers: handlers,
	}

	handlers.Register(b)

	return &tgBot, nil
}

func (b *TgBot) Start(ctx context.Context) error {
	logger.Info(ctx, "starting telegram bot...")
	b.bot.Start(ctx)
	return nil
}

func (b *TgBot) Stop(ctx context.Context) {
	logger.Info(ctx, "telegram bot stopped")
}
