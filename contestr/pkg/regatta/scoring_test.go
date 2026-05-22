package regatta

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestScoringSettingsUnmarshalBSON_defaultsMissingFields(t *testing.T) {
	data, err := bson.Marshal(bson.M{
		"solve_in_time_bonus": 100,
	})
	if err != nil {
		t.Fatalf("marshal scoring settings: %v", err)
	}

	var settings ScoringSettings
	if err := bson.Unmarshal(data, &settings); err != nil {
		t.Fatalf("unmarshal scoring settings: %v", err)
	}

	if settings.OvertakeBonus != DefaultOvertakeBonus {
		t.Fatalf("overtake bonus = %d, want default %d", settings.OvertakeBonus, DefaultOvertakeBonus)
	}
}

func TestScoringSettingsUnmarshalBSON_preservesExplicitZero(t *testing.T) {
	data, err := bson.Marshal(bson.M{
		"solve_in_time_bonus": 0,
		"overtake_bonus":      0,
	})
	if err != nil {
		t.Fatalf("marshal scoring settings: %v", err)
	}

	var settings ScoringSettings
	if err := bson.Unmarshal(data, &settings); err != nil {
		t.Fatalf("unmarshal scoring settings: %v", err)
	}

	if settings.SolveInTimeBonus != 0 || settings.OvertakeBonus != 0 {
		t.Fatalf("explicit zero bonuses were not preserved: %+v", settings)
	}
}

func TestScoringSettingsUnmarshalBSON_preservesExplicitFieldZero(t *testing.T) {
	data, err := bson.Marshal(bson.M{
		"solve_in_time_bonus": 0,
		"overtake_bonus":      100,
	})
	if err != nil {
		t.Fatalf("marshal scoring settings: %v", err)
	}

	var settings ScoringSettings
	if err := bson.Unmarshal(data, &settings); err != nil {
		t.Fatalf("unmarshal scoring settings: %v", err)
	}

	if settings.SolveInTimeBonus != 0 {
		t.Fatalf("solve in time bonus = %d, want explicit 0", settings.SolveInTimeBonus)
	}
}

func TestScoringSettingsUnmarshalBSON_ignoresLegacyFields(t *testing.T) {
	data, err := bson.Marshal(bson.M{
		"full_solve_bonus":     50,
		"solve_in_time_bonus":  80,
		"overtake_bonus":       10,
	})
	if err != nil {
		t.Fatalf("marshal scoring settings: %v", err)
	}

	var settings ScoringSettings
	if err := bson.Unmarshal(data, &settings); err != nil {
		t.Fatalf("unmarshal scoring settings: %v", err)
	}

	if settings.SolveInTimeBonus != 80 || settings.OvertakeBonus != 10 {
		t.Fatalf("settings = %+v, want solve 80 overtake 10", settings)
	}
}

func TestNormalizeScoringSettings_allZeroUsesDefaults(t *testing.T) {
	got := NormalizeScoringSettings(ScoringSettings{})
	want := DefaultScoringSettings()
	if got != want {
		t.Fatalf("got %+v, want defaults %+v", got, want)
	}
}
