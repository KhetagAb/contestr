package regatta

import (
	"context"
	"testing"
	"time"

	regattapkg "contestr/pkg/regatta"
)

type testTimetableRepo struct {
	timetable *regattapkg.ToursTimetable
}

func (r *testTimetableRepo) Create(_ context.Context, timetable *regattapkg.ToursTimetable) error {
	r.timetable = timetable
	return nil
}

func (r *testTimetableRepo) GetByContestID(_ context.Context, _ int) (*regattapkg.ToursTimetable, error) {
	return r.timetable, nil
}

func (r *testTimetableRepo) Update(_ context.Context, timetable *regattapkg.ToursTimetable) error {
	r.timetable = timetable
	return nil
}

func (r *testTimetableRepo) DeleteByContestID(_ context.Context, _ int) error {
	r.timetable = nil
	return nil
}

func TestUpdateActiveTourDurationAllowsFirstFactualTourBeforeContestStart(t *testing.T) {
	tourRepo := &startTourRepo{
		tours: []regattapkg.Tour{
			{
				ContestID:         42,
				Sequence:          1,
				Round:             1,
				DurationInSeconds: 1800,
			},
		},
	}
	contestRepo := &startContestRepo{
		contest: &regattapkg.Contest{
			ContestID: 42,
			StartTime: time.Now().Add(time.Hour),
		},
	}
	timetableRepo := &testTimetableRepo{
		timetable: &regattapkg.ToursTimetable{ContestId: 42},
	}

	service := NewRegatta(tourRepo, contestRepo, timetableRepo)
	view, err := service.UpdateActiveTourDuration(context.Background(), 42, 2400, TimetableViewOptions{})
	if err != nil {
		t.Fatalf("update active duration before start: %v", err)
	}

	if tourRepo.updatedSequence != 1 {
		t.Fatalf("updated sequence = %d, want 1", tourRepo.updatedSequence)
	}
	if tourRepo.updatedDuration != 2400 {
		t.Fatalf("updated duration = %d, want 2400", tourRepo.updatedDuration)
	}
	if len(view.TimelineSegments) != 1 || view.TimelineSegments[0].Duration != 2400 {
		t.Fatalf("view timeline = %+v, want one segment with duration 2400", view.TimelineSegments)
	}
}
