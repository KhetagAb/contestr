package regatta

import (
	"context"
	"errors"
	"time"

	"contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/mongo"
)

type TimetableViewOptions struct {
	ServerAutoStartAvailable bool
}

func (s *Regatta) GetTimetableView(ctx context.Context, contestID int, opts TimetableViewOptions) (*regatta.TimetableView, error) {
	timetable, err := s.LoadTimetable(ctx, contestID)
	if errors.Is(err, ErrTimetableNotFound) {
		return s.buildTimetableView(ctx, contestID, &regatta.ToursTimetable{
			ContestId: contestID,
			TourTimes: []regatta.TourConfig{},
		}, opts)
	}
	if err != nil {
		return nil, err
	}
	return s.buildTimetableView(ctx, contestID, timetable, opts)
}

func (s *Regatta) buildTimetableView(
	ctx context.Context,
	contestID int,
	timetable *regatta.ToursTimetable,
	opts TimetableViewOptions,
) (*regatta.TimetableView, error) {
	now := time.Now()

	view := &regatta.TimetableView{
		ContestId:          contestID,
		ServerNow:          now,
		TourTimes:          timetable.TourTimes,
		AutoStartAvailable: opts.ServerAutoStartAvailable,
		AutoStartEnabled:   opts.ServerAutoStartAvailable && timetable.AutoStartEnabled,
	}

	contest, err := s.contestRepo.GetByContestID(ctx, contestID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if contest != nil {
		view.ContestStartTime = &contest.StartTime
		elapsed := int(now.Sub(contest.StartTime).Seconds())
		if elapsed < 0 {
			elapsed = 0
		}
		view.ElapsedSeconds = elapsed
	}

	meta, next := regatta.BuildToursMeta(timetable.TourTimes, view.ElapsedSeconds)
	view.ToursMeta = meta
	view.NextTourNumber = next

	return view, nil
}

func (s *Regatta) SaveTimetable(
	ctx context.Context,
	contestID int,
	input regatta.SaveTimetableRequest,
	opts TimetableViewOptions,
) (*regatta.TimetableView, error) {
	if contestID <= 0 {
		return nil, ErrInvalidTimetable
	}
	if err := validateTourDurations(input.TourDurations); err != nil {
		return nil, err
	}

	existing, err := s.LoadTimetable(ctx, contestID)
	if err != nil && !errors.Is(err, ErrTimetableNotFound) {
		return nil, err
	}

	merged := regatta.ApplyTimetableSchedule(existing, input.TourDurations)
	timetable := regatta.ToursTimetable{
		ContestId:        contestID,
		TourTimes:        merged,
		AutoStartEnabled: input.AutoStartEnabled,
	}

	var saved *regatta.ToursTimetable
	if existing == nil {
		saved, err = s.insertTimetable(ctx, timetable)
	} else {
		saved, err = s.ReplaceTimetable(ctx, timetable)
	}
	if err != nil {
		return nil, err
	}

	return s.buildTimetableView(ctx, contestID, saved, opts)
}
