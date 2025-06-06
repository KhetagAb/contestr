package tgbot

import (
	"contestr/pkg/logger"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type MessageHandle struct{}

func NewMessageHandle() *MessageHandle {
	return &MessageHandle{}
}

func (h *MessageHandle) Register() (bot.HandlerType, string, bot.MatchType, bot.HandlerFunc) {
	return bot.HandlerTypeMessageText, "", bot.MatchTypeCommand, h.HandleMessage
}

func (h *MessageHandle) HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Я понимаю только команды. Используйте /help для получения списка команд.",
	})

	if err != nil {
		logger.Errorf(ctx, "failed to send message response: %v", err)
	}
}
