package wiresets

import (
	"contestr/internal/handlers"
	cfhandlers "contestr/internal/handlers/codeforces"
	regattahandlers "contestr/internal/handlers/regatta"
	"contestr/internal/handlers/tgbot"
	"contestr/internal/integrations/codeforces"
	"contestr/internal/integrations/ejudge"
	"contestr/internal/repository"
	"contestr/internal/services/regatta"
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

	repository.NewMongoClient,
	repository.NewMongoTourRepository,

	ejudge.NewContestXMLFetcher,
	wire.Bind(new(regatta.EjudgeParser), new(*ejudge.ContestXMLFetcher)),
	wire.Bind(new(regatta.TourRepository), new(*repository.MongoTourRepository)),
	regatta.NewRegatta,

	wire.Bind(new(regattahandlers.Regatta), new(*regatta.Regatta)),
	regattahandlers.NewContestHandle,
)
