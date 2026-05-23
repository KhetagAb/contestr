package timetable_sync

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"contestr/internal/repository"
	regattasvc "contestr/internal/services/regatta"
	"contestr/pkg/logger"
	regattapkg "contestr/pkg/regatta"
)

func TestMain(m *testing.M) {
	logger.Init("timetable-sync-test", "error")
	os.Exit(m.Run())
}

// mockRegistry возвращает фиксированный набор контестов.
type mockRegistry struct {
	contests map[string][]int
}

func (m *mockRegistry) GetSystem(_ int) (string, error)                    { return "test", nil }
func (m *mockRegistry) GetAllContests() map[string][]int                   { return m.contests }
func (m *mockRegistry) GetContest(_ int) (*repository.RegisteredContest, error) { return nil, nil }
func (m *mockRegistry) GetAllRegisteredContests() []repository.RegisteredContest { return nil }

// mockRegatta записывает вызовы AdvanceTimetable и LoadTimetable.
type mockRegatta struct {
	timetable    *regattapkg.ToursTimetable
	loadErr      error
	advancedIDs  []int
	advanceErr   error
}

func (m *mockRegatta) LoadTimetable(_ context.Context, contestID int) (*regattapkg.ToursTimetable, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.timetable, nil
}

func (m *mockRegatta) AdvanceTimetable(_ context.Context, contestID int, _ regattasvc.AdvanceMode, _ regattasvc.TimetableViewOptions) error {
	m.advancedIDs = append(m.advancedIDs, contestID)
	return m.advanceErr
}

func newTestService(reg *mockRegistry, rg *mockRegatta) *TimetableSyncService {
	return NewTimetableSyncService(reg, rg, time.Second, true)
}

func TestAdvanceDue_skipsAutoStartDisabled(t *testing.T) {
	reg := &mockRegistry{contests: map[string][]int{"test": {1}}}
	rg := &mockRegatta{
		timetable: &regattapkg.ToursTimetable{ContestId: 1, AutoStartEnabled: false},
	}
	svc := newTestService(reg, rg)
	svc.AdvanceDue(context.Background())

	if len(rg.advancedIDs) != 0 {
		t.Fatalf("AdvanceTimetable должен не вызываться при AutoStart=false, вызван для: %v", rg.advancedIDs)
	}
}

func TestAdvanceDue_skipsOnTimetableNotFound(t *testing.T) {
	reg := &mockRegistry{contests: map[string][]int{"test": {42}}}
	rg := &mockRegatta{loadErr: regattasvc.ErrTimetableNotFound}
	svc := newTestService(reg, rg)
	svc.AdvanceDue(context.Background())

	if len(rg.advancedIDs) != 0 {
		t.Fatalf("AdvanceTimetable не должен вызываться, если расписание не найдено")
	}
}

func TestAdvanceDue_advancesWhenAutoStartEnabled(t *testing.T) {
	reg := &mockRegistry{contests: map[string][]int{"test": {7}}}
	rg := &mockRegatta{
		timetable:  &regattapkg.ToursTimetable{ContestId: 7, AutoStartEnabled: true},
		advanceErr: regattasvc.ErrNothingToAdvance,
	}
	svc := newTestService(reg, rg)
	svc.AdvanceDue(context.Background())

	if len(rg.advancedIDs) != 1 || rg.advancedIDs[0] != 7 {
		t.Fatalf("AdvanceTimetable должен быть вызван для контеста 7, got: %v", rg.advancedIDs)
	}
}

func TestAdvanceDue_propagatesUnexpectedError(t *testing.T) {
	reg := &mockRegistry{contests: map[string][]int{"test": {3}}}
	unexpectedErr := errors.New("db connection lost")
	rg := &mockRegatta{loadErr: unexpectedErr}
	svc := newTestService(reg, rg)

	// AdvanceDue логирует ошибку, но не паникует
	svc.AdvanceDue(context.Background())
}
