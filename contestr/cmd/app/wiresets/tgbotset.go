package wiresets

import (
	"contestr/internal/handlers/tgbot"
	"contestr/internal/transport"
	"github.com/google/wire"
)

var TgBotSet = wire.NewSet(
	tgbot.NewStartHandle,
	tgbot.NewHelpHandle,
	tgbot.NewMessageHandle,
	tgbot.NewHandlers,
	transport.NewBot,
)
