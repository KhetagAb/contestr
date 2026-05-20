package regatta

import "go.mongodb.org/mongo-driver/bson"

const (
	ScoringModeBinary  = "binary"
	ScoringModePartial = "partial"

	BinaryOvertakeModeRetrospective = "retrospective"
	BinaryOvertakeModeDuringTour    = "during_tour_only"
)

const (
	DefaultFullSolveBonus   = 100
	DefaultSolveInTimeBonus = 100
	DefaultOvertakeBonus    = 100

	CurrentScoringSettingsVersion = 1
)

type ScoringSettings struct {
	Version            int    `bson:"version" json:"-"`
	Mode               string `bson:"mode" json:"mode"`
	BinaryOvertakeMode string `bson:"binary_overtake_mode" json:"binary_overtake_mode"`
	FullSolveBonus     int    `bson:"full_solve_bonus" json:"full_solve_bonus"`
	SolveInTimeBonus   int    `bson:"solve_in_time_bonus" json:"solve_in_time_bonus"`
	OvertakeBonus      int    `bson:"overtake_bonus" json:"overtake_bonus"`
}

func (settings *ScoringSettings) UnmarshalBSON(data []byte) error {
	type rawScoringSettings struct {
		Version            *int    `bson:"version"`
		Mode               *string `bson:"mode"`
		BinaryOvertakeMode *string `bson:"binary_overtake_mode"`
		FullSolveBonus     *int    `bson:"full_solve_bonus"`
		SolveInTimeBonus   *int    `bson:"solve_in_time_bonus"`
		OvertakeBonus      *int    `bson:"overtake_bonus"`
	}

	var raw rawScoringSettings
	if err := bson.Unmarshal(data, &raw); err != nil {
		return err
	}

	next := ScoringSettings{}
	if raw.Version != nil {
		next.Version = *raw.Version
	}
	if raw.Mode != nil {
		next.Mode = *raw.Mode
	}
	if raw.BinaryOvertakeMode != nil {
		next.BinaryOvertakeMode = *raw.BinaryOvertakeMode
	}
	if raw.FullSolveBonus != nil {
		next.FullSolveBonus = *raw.FullSolveBonus
	}
	if raw.SolveInTimeBonus != nil {
		next.SolveInTimeBonus = *raw.SolveInTimeBonus
	}
	if raw.OvertakeBonus != nil {
		next.OvertakeBonus = *raw.OvertakeBonus
	}

	*settings = NormalizeScoringSettings(next)
	return nil
}

func DefaultScoringSettings() ScoringSettings {
	return ScoringSettings{
		Version:            CurrentScoringSettingsVersion,
		Mode:               ScoringModeBinary,
		BinaryOvertakeMode: BinaryOvertakeModeRetrospective,
		FullSolveBonus:     DefaultFullSolveBonus,
		SolveInTimeBonus:   DefaultSolveInTimeBonus,
		OvertakeBonus:      DefaultOvertakeBonus,
	}
}

func NormalizeScoringSettings(settings ScoringSettings) ScoringSettings {
	defaults := DefaultScoringSettings()
	if settings.Mode == "" &&
		settings.BinaryOvertakeMode == "" &&
		settings.FullSolveBonus == 0 &&
		settings.SolveInTimeBonus == 0 &&
		settings.OvertakeBonus == 0 {
		return defaults
	}

	legacySettings := settings.Version == 0
	settings.Version = CurrentScoringSettingsVersion

	if settings.Mode != ScoringModePartial {
		settings.Mode = ScoringModeBinary
	}
	if settings.BinaryOvertakeMode != BinaryOvertakeModeDuringTour {
		settings.BinaryOvertakeMode = BinaryOvertakeModeRetrospective
	}
	if legacySettings && settings.FullSolveBonus == 0 {
		settings.FullSolveBonus = defaults.FullSolveBonus
	}
	if legacySettings && settings.SolveInTimeBonus == 0 {
		settings.SolveInTimeBonus = defaults.SolveInTimeBonus
	}
	if legacySettings && settings.OvertakeBonus == 0 {
		settings.OvertakeBonus = defaults.OvertakeBonus
	}
	if settings.FullSolveBonus < 0 {
		settings.FullSolveBonus = 0
	}
	if settings.SolveInTimeBonus < 0 {
		settings.SolveInTimeBonus = 0
	}
	if settings.OvertakeBonus < 0 {
		settings.OvertakeBonus = 0
	}

	return settings
}
