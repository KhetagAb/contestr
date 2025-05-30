package tgbot

import (
	"contestr/pkg/logger"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// HandleStart обрабатывает команду /start
func (h *Handlers) HandleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	userID := fmt.Sprint(update.Message.From.ID)

	welcomeMsg, err := handleStartCommand(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "error handling start command: %v", err)
		sendErrorMessage(ctx, b, chatID)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   welcomeMsg,
	})

	if err != nil {
		logger.Errorf(ctx, "error sending message: %v", err)
	}
}

func handleStartCommand(_ context.Context, _ string) (string, error) {
	return "Привет! Я бот для работы с Codeforces.\n" +
			"Используйте /help для получения списка команд.",
		nil
}
