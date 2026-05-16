package wiresets

import (
	"contestr/internal/auth"
	"contestr/internal/configs"
	"contestr/internal/handlers"
	adminhandlers "contestr/internal/handlers/admin"
	cfhandlers "contestr/internal/handlers/codeforces"
	regattahandlers "contestr/internal/handlers/regatta"
	"contestr/internal/handlers/tgbot"
	"contestr/internal/integrations"
	"contestr/internal/integrations/codeforces"
	"contestr/internal/integrations/ejudge"
	"contestr/internal/repository"
	"contestr/internal/services/contest_registry"
	"contestr/internal/services/contest_sync"
	"contestr/internal/services/regatta"
	"contestr/internal/transport"
	"contestr/pkg/config"
	"context"
	"time"

	"github.com/google/wire"
)

func NewContextProvider() context.Context {
	return context.Background()
}

func NewContestAdaptersMap(
	ejudgeAdapter *ejudge.EjudgeAdapter,
	codeforcesAdapter *codeforces.CodeforcesAdapter,
) map[string]integrations.ContestAdapter {
	return map[string]integrations.ContestAdapter{
		"ejudge":     ejudgeAdapter,
		"codeforces": codeforcesAdapter,
	}
}

func GetContestSyncInterval(cfg *configs.Config) time.Duration {
	return cfg.ContestSync.Interval
}

var All = wire.NewSet(
	NewContextProvider,
	config.NewConfig,

	wire.Bind(new(cfhandlers.Service), new(*codeforces.Service)),
	codeforces.NewService,

	cfhandlers.NewContestHandle,
	adminhandlers.NewLoginHandle,
	adminhandlers.NewMeHandle,
	handlers.NewHandlers,
	auth.NewService,
	transport.NewHTTPServer,

	tgbot.NewStartHandle,
	wire.Bind(new(tgbot.Regatta), new(*regatta.Regatta)),
	tgbot.NewRegattaStartTourHandle,
	wire.Bind(new(tgbot.ContestSyncService), new(*contest_sync.ContestSyncService)),
	tgbot.NewSyncContestsHandle,
	tgbot.NewHelpHandle,
	tgbot.NewMessageHandle,
	tgbot.NewHandlers,
	transport.NewBot,

	repository.NewMongoClient,
	repository.NewMongoTourRepository,
	repository.NewMongoContestRepository,
	repository.NewMongoCodeforcesHandleRepository,
	wire.Bind(new(regatta.TourRepository), new(*repository.MongoTourRepository)),
	wire.Bind(new(regatta.ContestRepository), new(*repository.MongoContestRepository)),
	wire.Bind(new(repository.ContestRepository), new(*repository.MongoContestRepository)),
	wire.Bind(new(repository.CodeforcesHandleRepository), new(*repository.MongoCodeforcesHandleRepository)),

	ejudge.NewContestXMLFetcher,
	ejudge.NewEjudgeAdapter,

	codeforces.NewCodeforcesAdapter,

	NewContestAdaptersMap,

	contest_registry.NewContestRegistry,

	GetContestSyncInterval,
	contest_sync.NewContestSyncService,

	regatta.NewRegatta,

	wire.Bind(new(regattahandlers.Regatta), new(*regatta.Regatta)),
	regattahandlers.NewContestHandle,
)
