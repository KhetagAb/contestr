package tgbot

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// HandleRegister обрабатывает команду /register <handle>
func (h *Handlers) HandleRegister(ctx context.Context, b *bot.Bot, update *models.Update) {
	//panic("implement me")
	//chatID := update.Message.Chat.ID
	//userID := update.Message.From.ID
	//username := update.Message.From.Username
	//
	//args := strings.Fields(update.Message.Text)
	//
	//if len(args) < 2 {
	//	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
	//		ChatID: chatID,
	//		Text:   "Пожалуйста, укажите ваш Codeforces handle. Пример: /register tourist",
	//	})
	//
	//	if err != nil {
	//		logger.Errorf(ctx, "failed to send error message: %v", err)
	//	}
	//	return
	//}
	//
	//handle := args[1]
	//
	//// Предполагаем, что у нас есть метод для обновления Codeforces handle
	//err := h.userUseCase.UpdateUserCodeforcesHandle(ctx, string(userID), handle)
	//
	//var responseMsg string
	//if err != nil {
	//	responseMsg = "Произошла ошибка при регистрации: " + err.Error()
	//	logger.Errorf(ctx, "failed to register user %s with handle %s: %v", username, handle, err)
	//} else {
	//	responseMsg = "Ваш Codeforces handle успешно зарегистрирован: " + handle
	//	logger.Infof(ctx, "user %s registered with handle %s", username, handle)
	//}
	//
	//_, err = b.SendMessage(ctx, &bot.SendMessageParams{
	//	ChatID: chatID,
	//	Text:   responseMsg,
	//})
	//
	//if err != nil {
	//	logger.Errorf(ctx, "failed to send registration response: %v", err)
	//}
}
