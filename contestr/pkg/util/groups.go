package util

import (
	"math/rand"
)

const (
	SwapBorderProbability = 0.4
)

func FormGroups(ratedParticipants []string, groupSize int) [][]string {
	n := len(ratedParticipants)
	if n == 0 {
		return [][]string{}
	}
	if groupSize <= 0 {
		groupSize = 1
	}

	participants := append([]string(nil), ratedParticipants...)

	for i := 0; i < n-1; i++ {
		if rand.Float64() <= SwapBorderProbability {
			participants[i], participants[i+1] = participants[i+1], participants[i]
		}
	}

	result := make([][]string, 0, n/groupSize+1)
	for i := 0; i < n; i += groupSize {
		intervalEnd := i + groupSize
		if intervalEnd > n {
			intervalEnd = n
		}
		result = append(result, participants[i:intervalEnd])
	}

	last := len(result) - 1
	if last >= 1 && len(result[last]) == 1 {
		prev := result[last-1]
		if len(prev) >= 3 {
			result[last] = []string{prev[2], result[last][0]}
			result[last-1] = prev[:2]
		} else {
			result[last-1] = append(append([]string(nil), prev...), result[last][0])
			result = result[:last]
		}
	}

	return result
}
