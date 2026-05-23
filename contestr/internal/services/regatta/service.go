package regatta

import (
	"context"
	"sync"

	"contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourRepository interface {
	Create(ctx context.Context, tour *regatta.Tour) (primitive.ObjectID, error)
	FindByContestID(ctx context.Context, contestID int) ([]regatta.Tour, error)
	UpdateDuration(ctx context.Context, contestID int, sequence int, durationSeconds int) error
}

type ContestRepository interface {
	GetByContestID(ctx context.Context, contestID int) (*regatta.Contest, error)
	GetParticipants(ctx context.Context, contestID int) (map[string]string, error)
	GetSubmissions(ctx context.Context, contestID int) ([]regatta.ContestSubmission, error)
}

type TimetableRepository interface {
	Create(ctx context.Context, timetable *regatta.ToursTimetable) error
	GetByContestID(ctx context.Context, contestID int) (*regatta.ToursTimetable, error)
	Update(ctx context.Context, timetable *regatta.ToursTimetable) error
	DeleteByContestID(ctx context.Context, contestID int) error
}

type Regatta struct {
	tourRepository      TourRepository
	contestRepo         ContestRepository
	timetableRepository TimetableRepository

	advanceLockMu sync.Mutex
	advanceLocks  map[int]*sync.Mutex
}

func NewRegatta(
	tourRepository TourRepository,
	contestRepo ContestRepository,
	timetableRepository TimetableRepository,
) *Regatta {
	return &Regatta{
		tourRepository:      tourRepository,
		contestRepo:         contestRepo,
		timetableRepository: timetableRepository,
	}
}
