package tgbot

import (
	"contestr/pkg/logger"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
	"strings"
	"time"
)

type Regatta interface {
	StartTour(ctx context.Context, contestId int, duration time.Duration) (string, error)
}

type RegattaStartTourHandle struct {
	regatta Regatta
}

func NewRegattaStartTourHandle(
	regatta Regatta,
) *RegattaStartTourHandle {
	return &RegattaStartTourHandle{
		regatta: regatta,
	}
}

func (h *RegattaStartTourHandle) Register() (bot.HandlerType, string, bot.MatchType, bot.HandlerFunc) {
	return bot.HandlerTypeMessageText, "start_tour", bot.MatchTypeCommand, h.Handle
}

func (h *RegattaStartTourHandle) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	message := "Неверный формат, необходимо: \n" +
		"/start_tour <contest_id> <duration_in_minutes>"
	parts := strings.Split(update.Message.Text, " ")[1:]

	if len(parts) == 2 {
		contestId, err1 := strconv.Atoi(parts[0])
		duration, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil {
			objectID, err := h.regatta.StartTour(ctx, contestId, time.Duration(duration)*time.Minute)
			if err != nil {
				message = fmt.Sprintf("Ошибка при добавлении нового тура: %v", err)
			} else {
				message = fmt.Sprintf("Начался тур... [ObjectId=%v]", objectID)
			}
		}
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})

	if err != nil {
		logger.Errorf(ctx, "error sending message: %v", err)
	}
}
