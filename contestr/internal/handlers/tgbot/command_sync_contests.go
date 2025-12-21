package tgbot

import (
	"contestr/internal/services/contest_sync"
	"contestr/pkg/logger"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strings"
)

type ContestSyncService interface {
	SyncAllContests(ctx context.Context) *contest_sync.SyncResult
}

type SyncContestsHandle struct {
	syncService ContestSyncService
}

func NewSyncContestsHandle(
	syncService ContestSyncService,
) *SyncContestsHandle {
	return &SyncContestsHandle{
		syncService: syncService,
	}
}

func (h *SyncContestsHandle) Register() (bot.HandlerType, string, bot.MatchType, bot.HandlerFunc) {
	return bot.HandlerTypeMessageText, "sync_contests", bot.MatchTypeCommand, h.Handle
}

func (h *SyncContestsHandle) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	logger.Infof(ctx, "manual sync requested by user %d", update.Message.From.ID)

	result := h.syncService.SyncAllContests(ctx)

	var message string
	if !result.HasContests {
		message = "Нет контестов для синхронизации. Проверьте конфигурацию."
	} else {
		var parts []string
		parts = append(parts, fmt.Sprintf("Синхронизация завершена:"))
		parts = append(parts, fmt.Sprintf("✅ Успешно: %d", result.SyncedCount))
		parts = append(parts, fmt.Sprintf("❌ Ошибок: %d", result.FailedCount))

		if len(result.ErrorMessages) > 0 {
			parts = append(parts, "")
			parts = append(parts, "Ошибки:")
			maxErrors := 5
			if len(result.ErrorMessages) < maxErrors {
				maxErrors = len(result.ErrorMessages)
			}
			for i := 0; i < maxErrors; i++ {
				parts = append(parts, fmt.Sprintf("• %s", result.ErrorMessages[i]))
			}
			if len(result.ErrorMessages) > maxErrors {
				parts = append(parts, fmt.Sprintf("... и еще %d ошибок", len(result.ErrorMessages)-maxErrors))
			}
		}

		message = strings.Join(parts, "\n")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})

	if err != nil {
		logger.Errorf(ctx, "error sending message: %v", err)
	}
}

