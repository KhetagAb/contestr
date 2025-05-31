package tgbot

import (
	"contestr/internal/configs"
	"contestr/internal/handlers/tgbot"
	"contestr/internal/transport"
	"contestr/internal/usecases"
	"fmt"
	_ "go.mongodb.org/mongo-driver/mongo"
)

func CreateBot(cfg *configs.Config) *transport.Bot {
	botUseCase := usecases.NewBotUseCase()

	botHandlers := tgbot.NewHandlers(botUseCase)
	bot, err := transport.NewBot(cfg, botHandlers)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize Telegram bot: %v", err))
	}
	bot.RegisterHandlers()
	return bot
}
