package regatta

import (
	"context"
	"slices"
	"testing"
	"time"

	regattapkg "contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type startTourRepo struct {
	tours           []regattapkg.Tour
	created         *regattapkg.Tour
	updatedSequence int
	updatedDuration int
}

func (r *startTourRepo) Create(_ context.Context, tour *regattapkg.Tour) (primitive.ObjectID, error) {
	created := *tour
	r.created = &created
	return primitive.NewObjectID(), nil
}

func (r *startTourRepo) FindByContestID(_ context.Context, _ int) ([]regattapkg.Tour, error) {
	return r.tours, nil
}

func (r *startTourRepo) UpdateDuration(_ context.Context, _ int, sequence int, durationSeconds int) error {
	r.updatedSequence = sequence
	r.updatedDuration = durationSeconds
	for i := range r.tours {
		if r.tours[i].Sequence == sequence {
			r.tours[i].DurationInSeconds = durationSeconds
		}
	}
	return nil
}

type startContestRepo struct {
	contest *regattapkg.Contest
}

func (r *startContestRepo) GetByContestID(_ context.Context, _ int) (*regattapkg.Contest, error) {
	return r.contest, nil
}

func (r *startContestRepo) GetParticipants(_ context.Context, _ int) (map[string]string, error) {
	result := make(map[string]string, len(r.contest.Participants))
	for _, participant := range r.contest.Participants {
		result[participant.ID] = participant.DisplayName
	}
	return result, nil
}

func (r *startContestRepo) GetSubmissions(_ context.Context, _ int) ([]regattapkg.ContestSubmission, error) {
	return r.contest.Submissions, nil
}

func TestStartTourUsesContestTourSettings(t *testing.T) {
	tourRepo := &startTourRepo{
		tours: []regattapkg.Tour{
			{
				ContestID:         42,
				Sequence:          1,
				Round:             1,
				DurationInSeconds: 1800,
				Problems:          []int{1, 2, 3},
				Groups:            map[string][]string{},
				GroupNumbers:      map[string]int{},
				GroupSize:         3,
			},
		},
	}
	contestRepo := &startContestRepo{
		contest: &regattapkg.Contest{
			ContestID:       42,
			ContestName:     "Configurable contest",
			System:          "codeforces",
			StartTime:       time.Now(),
			ScoringSettings: regattapkg.DefaultScoringSettings(),
			TourSettings: regattapkg.TourSettings{
				GroupSize:           2,
				ProblemsPerTour:     3,
				GroupShufflePercent: 0,
			},
			Participants: []regattapkg.ContestParticipant{
				{ID: "a", DisplayName: "A"},
				{ID: "b", DisplayName: "B"},
				{ID: "c", DisplayName: "C"},
				{ID: "d", DisplayName: "D"},
				{ID: "e", DisplayName: "E"},
				{ID: "f", DisplayName: "F"},
			},
		},
	}

	service := NewRegatta(tourRepo, contestRepo, nil)
	if _, err := service.StartTour(context.Background(), 42, 2400, StartTourOptions{}); err != nil {
		t.Fatalf("start tour: %v", err)
	}

	if tourRepo.created == nil {
		t.Fatal("expected tour to be created")
	}
	if tourRepo.created.GroupSize != 2 {
		t.Fatalf("group size = %d, want 2", tourRepo.created.GroupSize)
	}
	if !slices.Equal(tourRepo.created.Problems, []int{4, 5, 6}) {
		t.Fatalf("problems = %+v, want [4 5 6]", tourRepo.created.Problems)
	}
	for participant, group := range tourRepo.created.Groups {
		if len(group) != 2 {
			t.Fatalf("participant %s group size = %d, want 2: %+v", participant, len(group), group)
		}
	}
}

func TestParticipantsOrderedByRatingOrder_sortsByScore(t *testing.T) {
	tourRepo := &startTourRepo{
		tours: []regattapkg.Tour{
			{
				Sequence:  1,
				Round:     1,
				Problems:  []int{1},
				GroupSize: 2,
				Groups: map[string][]string{
					"p1": {"p1", "p2"}, "p2": {"p1", "p2"},
					"p3": {"p3", "p4"}, "p4": {"p3", "p4"},
				},
				GroupNumbers: map[string]int{"p1": 1, "p2": 1, "p3": 2, "p4": 2},
			},
		},
	}
	contestRepo := &startContestRepo{
		contest: &regattapkg.Contest{
			ContestID:       42,
			ScoringSettings: regattapkg.DefaultScoringSettings(),
			Participants: []regattapkg.ContestParticipant{
				{ID: "p1", DisplayName: "First"},
				{ID: "p2", DisplayName: "Second"},
				{ID: "p3", DisplayName: "Third"},
				{ID: "p4", DisplayName: "Fourth"},
			},
			Submissions: []regattapkg.ContestSubmission{
				{ParticipantID: "p1", ProblemID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 100},
				{ParticipantID: "p2", ProblemID: 1, Time: 200, Status: SubmissionStatusPartial, Points: 80},
				{ParticipantID: "p3", ProblemID: 1, Time: 300, Status: SubmissionStatusPartial, Points: 60},
				{ParticipantID: "p4", ProblemID: 1, Time: 400, Status: SubmissionStatusPartial, Points: 40},
			},
		},
	}
	service := NewRegatta(tourRepo, contestRepo, nil)
	participantsMap, _ := contestRepo.GetParticipants(context.Background(), 42)
	contest, _ := contestRepo.GetByContestID(context.Background(), 42)
	got, err := service.participantsOrderedByRating(
		context.Background(), 42, participantsMap, tourRepo.tours, contest,
	)
	if err != nil {
		t.Fatalf("participantsOrderedByRating: %v", err)
	}
	want := []Participant{"p1", "p2", "p3", "p4"}
	if !slices.Equal(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

func TestStartTourGroupsByRatingOrder(t *testing.T) {
	// After tour 1: p1=100, p2=80, p3=60, p4=40. With group_size=2 and no shuffle → (p1,p2) and (p3,p4).
	tourRepo := &startTourRepo{
		tours: []regattapkg.Tour{
			{
				ContestID:         42,
				Sequence:          1,
				Round:             1,
				DurationInSeconds: 0,
				Problems:          []int{1},
				GroupSize:         2,
				Groups: map[string][]string{
					"p1": {"p1", "p2"},
					"p2": {"p1", "p2"},
					"p3": {"p3", "p4"},
					"p4": {"p3", "p4"},
				},
				GroupNumbers: map[string]int{
					"p1": 1, "p2": 1, "p3": 2, "p4": 2,
				},
			},
		},
	}
	contestRepo := &startContestRepo{
		contest: &regattapkg.Contest{
			ContestID:       42,
			ContestName:     "Rating groups",
			System:          "codeforces",
			StartTime:       time.Now().Add(-2 * time.Hour),
			ScoringSettings: regattapkg.DefaultScoringSettings(),
			TourSettings: regattapkg.TourSettings{
				GroupSize:           2,
				ProblemsPerTour:     1,
				GroupShufflePercent: 0,
			},
			Participants: []regattapkg.ContestParticipant{
				{ID: "p1", DisplayName: "First"},
				{ID: "p2", DisplayName: "Second"},
				{ID: "p3", DisplayName: "Third"},
				{ID: "p4", DisplayName: "Fourth"},
			},
			Submissions: []regattapkg.ContestSubmission{
				{ParticipantID: "p1", ProblemID: 1, Time: 100, Status: SubmissionStatusPartial, Points: 100},
				{ParticipantID: "p2", ProblemID: 1, Time: 200, Status: SubmissionStatusPartial, Points: 80},
				{ParticipantID: "p3", ProblemID: 1, Time: 300, Status: SubmissionStatusPartial, Points: 60},
				{ParticipantID: "p4", ProblemID: 1, Time: 400, Status: SubmissionStatusPartial, Points: 40},
			},
		},
	}

	service := NewRegatta(tourRepo, contestRepo, nil)
	if _, err := service.StartTour(context.Background(), 42, 1000, StartTourOptions{}); err != nil {
		t.Fatalf("start tour: %v", err)
	}
	expectGroup := func(participant string, wantGroup int) {
		t.Helper()
		got := tourRepo.created.GroupNumbers[participant]
		if got != wantGroup {
			t.Fatalf("%s group = %d, want %d", participant, got, wantGroup)
		}
	}
	expectGroup("p1", 1)
	expectGroup("p2", 1)
	expectGroup("p3", 2)
	expectGroup("p4", 2)
}
