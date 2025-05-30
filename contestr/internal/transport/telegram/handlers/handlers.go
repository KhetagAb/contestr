package handlers

import (
	"contestr/pkg/logger"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

// HandleStart обрабатывает команду /start
func (h *Handlers) HandleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username

	logger.Infof(ctx, "user %s started the bot", username)

	welcomeMsg := "Привет! Я бот для работы с Codeforces.\n" +
		"Используйте /help для получения списка команд."
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   welcomeMsg,
	})

	if err != nil {
		logger.Errorf(ctx, "failed to send welcome message: %v", err)
	}
}

// HandleHelp обрабатывает команду /help
func (h *Handlers) HandleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	helpMsg := "Доступные команды:\n" +
		"/start - Начать работу с ботом\n" +
		"/help - Показать это сообщение\n" +
		"/register <handle> - Зарегистрировать Codeforces аккаунт\n"

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   helpMsg,
	})

	if err != nil {
		logger.Errorf(ctx, "failed to send help message: %v", err)
	}
}

// HandleRegister обрабатывает команду /register <handle>
func (h *Handlers) HandleRegister(ctx context.Context, b *bot.Bot, update *models.Update) {
	panic("implement me")
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

func (h *Handlers) HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func (h *Handlers) HandleUnknownCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Неизвестная команда. Используйте /help для получения списка команд.",
	})

	if err != nil {
		logger.Errorf(ctx, "failed to send unknown command response: %v", err)
	}
}
