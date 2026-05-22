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
	"contestr/internal/handlers/contests"
	"contestr/internal/services/contest_admin"
	"contestr/internal/services/contest_registry"
	"contestr/internal/services/contest_sync"
	"contestr/internal/services/regatta"
	"contestr/internal/services/timetable_sync"
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

func GetTimetableSyncInterval(cfg *configs.Config) timetable_sync.Interval {
	return timetable_sync.Interval(cfg.TimetableSync.Interval)
}

func NewTimetableSyncServiceProvider(
	registry contest_registry.ContestRegistry,
	contestRepo repository.ContestRepository,
	regattaService *regatta.Regatta,
	cfg *configs.Config,
) *timetable_sync.TimetableSyncService {
	return timetable_sync.NewTimetableSyncService(
		registry,
		contestRepo,
		regattaService,
		GetTimetableSyncInterval(cfg),
		cfg.TimetableSync.Interval > 0,
	)
}

var All = wire.NewSet(
	NewContextProvider,
	config.NewConfig,

	wire.Bind(new(cfhandlers.Service), new(*codeforces.Service)),
	codeforces.NewService,

	cfhandlers.NewContestHandle,
	adminhandlers.NewLoginHandle,
	adminhandlers.NewMeHandle,
	wire.Bind(new(adminhandlers.TimetableService), new(*regatta.Regatta)),
	adminhandlers.NewTimetableHandle,
	adminhandlers.NewContestsHandle,
	contests.NewListHandle,
	wire.Bind(new(contest_admin.TourDeleter), new(*repository.MongoTourRepository)),
	wire.Bind(new(contest_admin.TimetableDeleter), new(*repository.TourTimetableRepository)),
	contest_admin.NewService,
	handlers.NewHandlers,
	auth.NewService,
	transport.NewHTTPServer,

	tgbot.NewStartHandle,
	wire.Bind(new(tgbot.ContestSyncService), new(*contest_sync.ContestSyncService)),
	tgbot.NewSyncContestsHandle,
	tgbot.NewHelpHandle,
	tgbot.NewMessageHandle,
	tgbot.NewHandlers,
	transport.NewBot,

	repository.NewMongoClient,
	repository.NewMongoTourRepository,
	repository.NewTourTimetableRepository,
	repository.NewMongoContestRepository,
	repository.NewMongoRegisteredContestRepository,
	repository.NewMongoCodeforcesHandleRepository,
	wire.Bind(new(repository.RegisteredContestRepository), new(*repository.MongoRegisteredContestRepository)),
	wire.Bind(new(regatta.TourRepository), new(*repository.MongoTourRepository)),
	wire.Bind(new(regatta.TimetableRepository), new(*repository.TourTimetableRepository)),
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
	GetTimetableSyncInterval,
	wire.Bind(new(timetable_sync.RegattaService), new(*regatta.Regatta)),
	NewTimetableSyncServiceProvider,

	regatta.NewRegatta,

	wire.Bind(new(regattahandlers.Regatta), new(*regatta.Regatta)),
	regattahandlers.NewContestHandle,
	wire.Bind(new(regattahandlers.TimetableViewer), new(*regatta.Regatta)),
	regattahandlers.NewTimetableHandle,
	wire.Bind(new(regattahandlers.HandleLister), new(*repository.MongoCodeforcesHandleRepository)),
	regattahandlers.NewParticipantsHandle,
)
