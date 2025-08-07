package wiresets

import (
	"contestr/internal/handlers"
	cfhandlers "contestr/internal/handlers/codeforces"
	"contestr/internal/handlers/regatta"
	"contestr/internal/handlers/tgbot"
	"contestr/internal/integrations/codeforces"
	"contestr/internal/integrations/ejudge"
	"contestr/internal/transport"
	"contestr/pkg/config"
	"context"
	"github.com/google/wire"
)

func NewContextProvider() context.Context {
	return context.Background()
}

var All = wire.NewSet(
	NewContextProvider,
	config.NewConfig,

	wire.Bind(new(cfhandlers.Service), new(*codeforces.Service)),
	codeforces.NewService,

	cfhandlers.NewContestHandle,
	handlers.NewHandlers,
	handlers.NewHelloHandle,
	transport.NewHTTPServer,

	tgbot.NewStartHandle,
	tgbot.NewHelpHandle,
	tgbot.NewMessageHandle,
	tgbot.NewHandlers,
	transport.NewBot,

	ejudge.NewContestXMLFetcher,
	regatta.NewContestHandle,
)
