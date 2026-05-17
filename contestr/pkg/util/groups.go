package util

import (
	"math/rand"
)

const (
	SwapBorderProbability = 0.4
)

func FormGroups(ratedParticipants []string, groupSize int) [][]string {
	n := len(ratedParticipants)
	result := make([][]string, 0, n/groupSize+1)

	for i := 0; i < n-1; i++ {
		if rand.Float64() <= SwapBorderProbability {
			ratedParticipants[i], ratedParticipants[i+1] = ratedParticipants[i+1], ratedParticipants[i]
		}
	}

	for i := 0; i < n; i += groupSize {
		intervalEnd := i + groupSize
		if intervalEnd > n {
			intervalEnd = n
		}

		result = append(result, ratedParticipants[i:intervalEnd])
	}

	if len(result[len(result)-1]) == 1 {
		result[len(result)-1] = []string{
			result[len(result)-2][2], result[len(result)-1][0],
		}
		result[len(result)-2] = []string{
			result[len(result)-2][0], result[len(result)-2][1],
		}
	}

	return result
}
