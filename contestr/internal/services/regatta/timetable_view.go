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
			ContestId:    contestID,
			PendingSlots: []regatta.ScheduleSlot{},
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
		PendingSlots:       timetable.PendingSlots,
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

	tours, err := s.loadToursSorted(ctx, contestID)
	if err != nil {
		return nil, err
	}

	view.TimelineSegments = regatta.BuildTimelineSegments(tours, timetable.PendingSlots, view.ElapsedSeconds)
	view.NextTourNumber = regatta.NextCompetitiveRound(tours, timetable.PendingSlots)

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
	if err := validatePendingSlots(input.PendingSlots); err != nil {
		return nil, err
	}

	existing, err := s.LoadTimetable(ctx, contestID)
	if err != nil && !errors.Is(err, ErrTimetableNotFound) {
		return nil, err
	}

	timetable := regatta.ToursTimetable{
		ContestId:        contestID,
		PendingSlots:     input.PendingSlots,
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
