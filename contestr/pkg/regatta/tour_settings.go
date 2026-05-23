package regatta

const (
	DefaultGroupSize           = 3
	DefaultProblemsPerTour     = 2
	DefaultGroupShufflePercent = 20
)

type TourSettings struct {
	GroupSize           int `bson:"group_size" json:"group_size"`
	ProblemsPerTour     int `bson:"problems_per_tour" json:"problems_per_tour"`
	GroupShufflePercent int `bson:"group_shuffle_percent" json:"group_shuffle_percent"`
}

func DefaultTourSettings() TourSettings {
	return TourSettings{
		GroupSize:           DefaultGroupSize,
		ProblemsPerTour:     DefaultProblemsPerTour,
		GroupShufflePercent: DefaultGroupShufflePercent,
	}
}

func (settings TourSettings) Normalize() TourSettings {
	defaults := DefaultTourSettings()
	if settings.GroupSize <= 0 {
		settings.GroupSize = defaults.GroupSize
	}
	if settings.ProblemsPerTour <= 0 {
		settings.ProblemsPerTour = defaults.ProblemsPerTour
	}
	settings.GroupShufflePercent = max(0, min(100, settings.GroupShufflePercent))
	return settings
}

// NormalizeTourSettings — функциональная форма для обратной совместимости.
func NormalizeTourSettings(settings TourSettings) TourSettings {
	return settings.Normalize()
}

func (settings TourSettings) GroupShuffleProbability() float64 {
	return float64(settings.Normalize().GroupShufflePercent) / 100
}
