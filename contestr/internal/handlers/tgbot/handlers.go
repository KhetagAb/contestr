package tgbot

import (
	"contestr/pkg/logger"
	"context"
	"github.com/go-telegram/bot"
)

type Handler interface {
	Register() (bot.HandlerType, string, bot.MatchType, bot.HandlerFunc)
}

type Handlers struct {
	handlers []Handler
}

func (h *Handlers) Register(b *bot.Bot) {
	for _, handler := range h.handlers {
		b.RegisterHandler(handler.Register())
	}
}

func NewHandlers(
	start *StartHandle,
	help *HelpHandle,
	message *MessageHandle,
) *Handlers {
	return &Handlers{
		handlers: []Handler{start, help, message},
	}
}

func sendErrorMessage(ctx context.Context, b *bot.Bot, chatID int64) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Произошла ошибка при обработке команды. Попробуйте позже.",
	})
	if err != nil {
		logger.Errorf(ctx, "error sending error message: %v", err)
	}
}
