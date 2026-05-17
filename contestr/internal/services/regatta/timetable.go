package regatta

import (
	"context"
	"errors"
	"fmt"

	"contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrTimetableNotFound      = errors.New("расписание не найдено")
	ErrTimetableAlreadyExists = errors.New("расписание уже существует")
	ErrTourNotFound           = errors.New("тур не найден в расписании")
	ErrInvalidTimetable       = errors.New("некорректное расписание")
)

type TimetableRepository interface {
	Create(ctx context.Context, timetable *regatta.ToursTimetable) error
	GetByContestID(ctx context.Context, contestID int) (*regatta.ToursTimetable, error)
	Update(ctx context.Context, timetable *regatta.ToursTimetable) error
	DeleteByContestID(ctx context.Context, contestID int) error
}

func (s *Regatta) CreateTimetable(ctx context.Context, timetable regatta.ToursTimetable) (*regatta.ToursTimetable, error) {
	if err := validateTimetable(timetable); err != nil {
		return nil, err
	}

	if _, err := s.timetableRepository.GetByContestID(ctx, timetable.ContestId); err == nil {
		return nil, ErrTimetableAlreadyExists
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("не удалось проверить расписание: %w", err)
	}

	if err := s.timetableRepository.Create(ctx, &timetable); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrTimetableAlreadyExists
		}
		return nil, fmt.Errorf("не удалось создать расписание: %w", err)
	}

	return &timetable, nil
}

func (s *Regatta) GetTimetable(ctx context.Context, contestID int) (*regatta.ToursTimetable, error) {
	timetable, err := s.timetableRepository.GetByContestID(ctx, contestID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTimetableNotFound
		}
		return nil, fmt.Errorf("не удалось получить расписание для контеста %d: %w", contestID, err)
	}
	return timetable, nil
}

func (s *Regatta) UpdateTimetable(ctx context.Context, timetable regatta.ToursTimetable) (*regatta.ToursTimetable, error) {
	if err := validateTimetable(timetable); err != nil {
		return nil, err
	}

	if _, err := s.GetTimetable(ctx, timetable.ContestId); err != nil {
		return nil, err
	}

	if err := s.timetableRepository.Update(ctx, &timetable); err != nil {
		return nil, fmt.Errorf("не удалось обновить расписание: %w", err)
	}

	return &timetable, nil
}

func (s *Regatta) DeleteTimetable(ctx context.Context, contestID int) error {
	if contestID <= 0 {
		return fmt.Errorf("%w: contest_id должен быть положительным", ErrInvalidTimetable)
	}

	if err := s.timetableRepository.DeleteByContestID(ctx, contestID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrTimetableNotFound
		}
		return fmt.Errorf("не удалось удалить расписание: %w", err)
	}

	return nil
}

func (s *Regatta) MoveTimetableTour(ctx context.Context, contestID int, tourNumber int, newStartTime int) (*regatta.ToursTimetable, error) {
	if contestID <= 0 {
		return nil, fmt.Errorf("%w: contest_id должен быть положительным", ErrInvalidTimetable)
	}
	if tourNumber <= 0 {
		return nil, fmt.Errorf("%w: tour_number должен быть положительным", ErrInvalidTimetable)
	}
	if newStartTime < 0 {
		return nil, fmt.Errorf("%w: start_time должен быть неотрицательным", ErrInvalidTimetable)
	}

	timetable, err := s.GetTimetable(ctx, contestID)
	if err != nil {
		return nil, err
	}

	position := tourNumber - 1
	if position >= len(timetable.TourTimes) {
		return nil, ErrTourNotFound
	}

	delta := newStartTime - timetable.TourTimes[position].StartTime
	if delta < 0 {
		return nil, fmt.Errorf("%w: новое время начала должно быть больше старого", ErrInvalidTimetable)
	}
	for i := position; i < len(timetable.TourTimes); i++ {
		nextStartTime := timetable.TourTimes[i].StartTime + delta
		if nextStartTime < 0 {
			return nil, fmt.Errorf("%w: перенос тура сделает start_time отрицательным", ErrInvalidTimetable)
		}
		timetable.TourTimes[i].StartTime = nextStartTime
	}

	if err := validateTimetable(*timetable); err != nil {
		return nil, err
	}

	if err := s.timetableRepository.Update(ctx, timetable); err != nil {
		return nil, fmt.Errorf("не удалось обновить расписание после переноса тура: %w", err)
	}

	return timetable, nil
}

func (s *Regatta) GetFirstNotStartedTimetableTour(ctx context.Context, contestID int) (*regatta.TourConfig, error) {
	timetable, err := s.GetTimetable(ctx, contestID)
	if err != nil {
		return nil, err
	}

	_, tour, ok := timetable.FirstNotStartedTour()
	if !ok {
		return nil, ErrTourNotFound
	}

	return &tour, nil
}

func validateTimetable(timetable regatta.ToursTimetable) error {
	if timetable.ContestId <= 0 {
		return fmt.Errorf("%w: contest_id должен быть положительным", ErrInvalidTimetable)
	}

	for i, tour := range timetable.TourTimes {
		if tour.StartTime < 0 {
			return fmt.Errorf("%w: у тура %d start_time должен быть неотрицательным", ErrInvalidTimetable, i+1)
		}
		if tour.Duration <= 0 {
			return fmt.Errorf("%w: у тура %d duration должен быть положительным", ErrInvalidTimetable, i+1)
		}
		if i > 0 {
			prevTour := timetable.TourTimes[i-1]
			if prevTour.StartTime+prevTour.Duration > tour.StartTime {
				return fmt.Errorf("%w: туры %d и %d пересекаются по времени", ErrInvalidTimetable, i, i+1)
			}
		}
	}

	return nil
}
