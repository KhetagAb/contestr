package delivery

import (
	"contestr/pkg/logger"
	"context"

	"contestr/internal/config"
	"contestr/internal/handlers/tgbot"
	"github.com/go-telegram/bot"
)

type Bot struct {
	bot      *bot.Bot
	cfg      *config.Config
	handlers *tgbot.Handlers
}

func NewBot(cfg *config.Config, handlers *tgbot.Handlers) (*Bot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.HandleUnknownCommand),
		bot.WithDebug(),
	}

	b, err := bot.New(cfg.Telegram.Token, opts...)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:      b,
		cfg:      cfg,
		handlers: handlers,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	b.registerHandlers()

	logger.Info(ctx, "starting Telegram bot...")
	b.bot.Start(ctx)
	return nil
}

func (b *Bot) Stop(ctx context.Context) {
	logger.Info(ctx, "telegram bot stopped")
}

func (b *Bot) registerHandlers() {
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, b.handlers.HandleStart)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, b.handlers.HandleHelp)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "register", bot.MatchTypeCommand, b.handlers.HandleRegister)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeCommand, b.handlers.HandleMessage)
}
