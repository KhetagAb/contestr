package tgbot

import (
	"contestr/pkg/logger"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// HandleHelp обрабатывает команду /help
func (h *Handlers) HandleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	userID := fmt.Sprint(update.Message.From.ID)

	msg, err := handleHelpCommand(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "error handling start command: %v", err)
		sendErrorMessage(ctx, b, chatID)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})

	if err != nil {
		logger.Errorf(ctx, "error sending message: %v", err)
	}
}

func handleHelpCommand(_ context.Context, _ string) (string, error) {
	return "Доступные команды:\n" +
			"/start - Начать работу с ботом\n" +
			"/help - Показать это сообщение\n" +
			"/register <handle> - Зарегистрировать Codeforces аккаунт" +
			"",
		nil
}
