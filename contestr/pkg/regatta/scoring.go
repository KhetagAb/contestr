package regatta

import "go.mongodb.org/mongo-driver/bson"

const (
	DefaultSolveInTimeBonus = 100
	DefaultOvertakeBonus    = 100
)

type ScoringSettings struct {
	SolveInTimeBonus int `bson:"solve_in_time_bonus" json:"solve_in_time_bonus"`
	OvertakeBonus    int `bson:"overtake_bonus" json:"overtake_bonus"`
}

func (settings *ScoringSettings) UnmarshalBSON(data []byte) error {
	type rawScoringSettings struct {
		SolveInTimeBonus *int `bson:"solve_in_time_bonus"`
		OvertakeBonus    *int `bson:"overtake_bonus"`
	}

	var raw rawScoringSettings
	if err := bson.Unmarshal(data, &raw); err != nil {
		return err
	}

	*settings = normalizeScoringFromBSON(raw.SolveInTimeBonus, raw.OvertakeBonus)
	return nil
}

func DefaultScoringSettings() ScoringSettings {
	return ScoringSettings{
		SolveInTimeBonus: DefaultSolveInTimeBonus,
		OvertakeBonus:    DefaultOvertakeBonus,
	}
}

func NormalizeScoringSettings(settings ScoringSettings) ScoringSettings {
	if settings.SolveInTimeBonus == 0 && settings.OvertakeBonus == 0 {
		return DefaultScoringSettings()
	}

	if settings.SolveInTimeBonus < 0 {
		settings.SolveInTimeBonus = 0
	}
	if settings.OvertakeBonus < 0 {
		settings.OvertakeBonus = 0
	}

	return settings
}

func normalizeScoringFromBSON(solveInTime, overtake *int) ScoringSettings {
	defaults := DefaultScoringSettings()
	if solveInTime == nil && overtake == nil {
		return defaults
	}

	next := defaults
	if solveInTime != nil {
		next.SolveInTimeBonus = nonNegativeBonus(*solveInTime)
	}
	if overtake != nil {
		next.OvertakeBonus = nonNegativeBonus(*overtake)
	}
	return next
}

func nonNegativeBonus(value int) int {
	if value < 0 {
		return 0
	}
	return value
}
