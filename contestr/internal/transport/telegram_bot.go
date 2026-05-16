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

func NewBot(ctx context.Context, cfg *configs.Config, handlers *tgbot.Handlers) (*TgBot, error) {
	if !cfg.Telegram.Enabled() {
		logger.Info(ctx, "telegram bot is disabled: token is not configured")
		return &TgBot{cfg: cfg, handlers: handlers}, nil
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandle),
		bot.WithDebug(),
	}

	b, err := bot.New(cfg.Telegram.Token, opts...)
	if err != nil {
		logger.Warnf(ctx, "telegram bot is disabled: failed to connect: %v", err)
		return &TgBot{cfg: cfg, handlers: handlers}, nil
	}

	tgBot := TgBot{
		bot:      b,
		cfg:      cfg,
		handlers: handlers,
	}

	handlers.Register(b)

	logger.Info(ctx, "telegram bot initialized")
	return &tgBot, nil
}

func (b *TgBot) Enabled() bool {
	return b != nil && b.bot != nil
}

func (b *TgBot) Start(ctx context.Context) error {
	if !b.Enabled() {
		return nil
	}

	logger.Info(ctx, "starting telegram bot...")
	b.bot.Start(ctx)
	return nil
}

func (b *TgBot) Stop(ctx context.Context) {
	if !b.Enabled() {
		return
	}

	logger.Info(ctx, "telegram bot stopped")
}
