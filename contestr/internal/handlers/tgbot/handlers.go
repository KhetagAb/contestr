package tgbot

import (
	"contestr/internal/usecases"
	"contestr/pkg/logger"
	"context"
	"github.com/go-telegram/bot"
)

type Handlers struct {
	botUseCase *usecases.BotUseCase
}

func NewHandlers(botUseCase *usecases.BotUseCase) *Handlers {
	return &Handlers{
		botUseCase: botUseCase,
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
