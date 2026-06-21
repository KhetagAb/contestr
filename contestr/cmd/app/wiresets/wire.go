package wiresets

import (
	"contestr/internal/auth"
	"contestr/internal/configs"
	"contestr/internal/handlers"
	adminhandlers "contestr/internal/handlers/admin"
	regattahandlers "contestr/internal/handlers/regatta"
	"contestr/internal/integrations"
	"contestr/internal/integrations/codeforces"
	"contestr/internal/integrations/ejudge"
	"contestr/internal/repository"
	"contestr/internal/handlers/contests"
	"contestr/internal/services/contest_admin"
	"contestr/internal/services/contest_registry"
	"contestr/internal/services/contest_sync"
	"contestr/internal/services/problem_statement"
	"contestr/internal/services/regatta"
	"contestr/internal/storage/objectstorage"
	"contestr/internal/services/timetable_sync"
	"contestr/internal/transport"
	"contestr/pkg/config"
	"context"

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

func NewContestSyncServiceProvider(
	registry contest_registry.ContestRegistry,
	adapters map[string]integrations.ContestAdapter,
	contestRepo repository.ContestRepository,
	cfg *configs.Config,
) *contest_sync.ContestSyncService {
	return contest_sync.NewContestSyncService(
		registry,
		adapters,
		contestRepo,
		cfg.ContestSync.Interval,
		cfg.ContestSync.IntervalBeforeStart,
	)
}

func NewTimetableSyncServiceProvider(
	registry contest_registry.ContestRegistry,
	regattaService *regatta.Regatta,
	cfg *configs.Config,
) *timetable_sync.TimetableSyncService {
	return timetable_sync.NewTimetableSyncService(
		registry,
		regattaService,
		cfg.TimetableSync.Interval,
		cfg.TimetableSync.Interval > 0,
	)
}

var All = wire.NewSet(
	NewContextProvider,
	config.NewConfig,

	codeforces.NewService,

	adminhandlers.NewLoginHandle,
	adminhandlers.NewMeHandle,
	wire.Bind(new(adminhandlers.TimetableService), new(*regatta.Regatta)),
	adminhandlers.NewTimetableHandle,
	adminhandlers.NewContestsHandle,
	adminhandlers.NewProblemStatementsHandle,
	contests.NewListHandle,
	wire.Bind(new(contest_admin.TourDeleter), new(*repository.MongoTourRepository)),
	wire.Bind(new(contest_admin.TimetableDeleter), new(*repository.TourTimetableRepository)),
	wire.Bind(new(contest_admin.ProblemStatementDeleter), new(*problem_statement.Service)),
	wire.Bind(new(problem_statement.BlobStore), new(*objectstorage.Client)),
	objectstorage.NewClientOptional,
	repository.NewMongoProblemStatementRepository,
	wire.Bind(new(repository.ProblemStatementRepository), new(*repository.MongoProblemStatementRepository)),
	wire.Bind(new(problem_statement.TourLister), new(*repository.MongoTourRepository)),
	problem_statement.NewService,
	contest_admin.NewService,
	handlers.NewHandlers,
	auth.NewService,
	transport.NewHTTPServer,

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

	NewContestSyncServiceProvider,
	wire.Bind(new(timetable_sync.RegattaService), new(*regatta.Regatta)),
	NewTimetableSyncServiceProvider,

	regatta.NewRegatta,

	wire.Bind(new(regattahandlers.Regatta), new(*regatta.Regatta)),
	regattahandlers.NewContestHandle,
	wire.Bind(new(regattahandlers.TimetableViewer), new(*regatta.Regatta)),
	regattahandlers.NewTimetableHandle,
	wire.Bind(new(regattahandlers.HandleLister), new(*repository.MongoCodeforcesHandleRepository)),
	regattahandlers.NewParticipantsHandle,
	regattahandlers.NewProblemStatementsHandle,
)
