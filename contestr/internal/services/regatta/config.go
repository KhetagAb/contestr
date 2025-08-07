package regatta

import (
	"contestr/pkg/regatta"
)

func EmptyResultFromTour(tc regatta.Tour) TourResult {
	return TourResult{
		Tour:            tc,
		Results:         make(map[Participant]ContestResult),
		ProblemsMapping: tc.ProblemsIDsToNameMapping(tc.Problems),
	}
}

//	Tour: Tour{
//		Name:      tc.Name,
//		Index:     tc.Index,
//		StartTime: util.ParseTimeOrPanic(tc.StartTime),
//		Duration:  time.Duration(tc.Duration) * time.Minute,
//		Groups:    ConvertGroups(tc.Groups),
//		GroupSize: func() int {
//			if len(tc.Groups) == 0 {
//				return 0
//			}
//			return len(tc.Groups[0])
//		}(),
//		Problems:  tc.Problems,
//		ContestID: tc.ContestID,
//	},
func ConvertGroups(groups [][]int) map[Participant]Group {
	result := make(map[Participant]Group)

	for _, group := range groups {
		for _, participantID := range group {
			result[participantID] = group
		}
	}

	return result
}
