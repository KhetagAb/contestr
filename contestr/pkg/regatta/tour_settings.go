package regatta

const (
	DefaultGroupSize           = 3
	DefaultProblemsPerTour     = 2
	DefaultGroupShufflePercent = 40
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

func NormalizeTourSettings(settings TourSettings) TourSettings {
	defaults := DefaultTourSettings()
	if settings.GroupSize <= 0 {
		settings.GroupSize = defaults.GroupSize
	}
	if settings.ProblemsPerTour <= 0 {
		settings.ProblemsPerTour = defaults.ProblemsPerTour
	}
	if settings.GroupShufflePercent < 0 {
		settings.GroupShufflePercent = 0
	}
	if settings.GroupShufflePercent > 100 {
		settings.GroupShufflePercent = 100
	}
	return settings
}

func (settings TourSettings) GroupShuffleProbability() float64 {
	normalized := NormalizeTourSettings(settings)
	return float64(normalized.GroupShufflePercent) / 100
}
