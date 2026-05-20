package regatta

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestScoringSettingsUnmarshalBSON_defaultsMissingFields(t *testing.T) {
	data, err := bson.Marshal(bson.M{
		"mode":             ScoringModePartial,
		"full_solve_bonus": 100,
	})
	if err != nil {
		t.Fatalf("marshal scoring settings: %v", err)
	}

	var settings ScoringSettings
	if err := bson.Unmarshal(data, &settings); err != nil {
		t.Fatalf("unmarshal scoring settings: %v", err)
	}

	if settings.Mode != ScoringModePartial {
		t.Fatalf("mode = %q, want %q", settings.Mode, ScoringModePartial)
	}
	if settings.SolveInTimeBonus != DefaultSolveInTimeBonus {
		t.Fatalf("solve in time bonus = %d, want default %d", settings.SolveInTimeBonus, DefaultSolveInTimeBonus)
	}
	if settings.OvertakeBonus != DefaultOvertakeBonus {
		t.Fatalf("overtake bonus = %d, want default %d", settings.OvertakeBonus, DefaultOvertakeBonus)
	}
}

func TestScoringSettingsUnmarshalBSON_preservesExplicitZero(t *testing.T) {
	data, err := bson.Marshal(bson.M{
		"version":              CurrentScoringSettingsVersion,
		"mode":                 ScoringModeBinary,
		"binary_overtake_mode": BinaryOvertakeModeRetrospective,
		"full_solve_bonus":     0,
		"solve_in_time_bonus":  0,
		"overtake_bonus":       0,
	})
	if err != nil {
		t.Fatalf("marshal scoring settings: %v", err)
	}

	var settings ScoringSettings
	if err := bson.Unmarshal(data, &settings); err != nil {
		t.Fatalf("unmarshal scoring settings: %v", err)
	}

	if settings.FullSolveBonus != 0 || settings.SolveInTimeBonus != 0 || settings.OvertakeBonus != 0 {
		t.Fatalf("explicit zero bonuses were not preserved: %+v", settings)
	}
}

func TestScoringSettingsUnmarshalBSON_defaultsLegacyZeroBonuses(t *testing.T) {
	data, err := bson.Marshal(bson.M{
		"mode":                 ScoringModeBinary,
		"binary_overtake_mode": BinaryOvertakeModeRetrospective,
		"full_solve_bonus":     100,
		"solve_in_time_bonus":  0,
		"overtake_bonus":       100,
	})
	if err != nil {
		t.Fatalf("marshal scoring settings: %v", err)
	}

	var settings ScoringSettings
	if err := bson.Unmarshal(data, &settings); err != nil {
		t.Fatalf("unmarshal scoring settings: %v", err)
	}

	if settings.SolveInTimeBonus != DefaultSolveInTimeBonus {
		t.Fatalf("legacy solve in time bonus = %d, want default %d", settings.SolveInTimeBonus, DefaultSolveInTimeBonus)
	}
}
