package util

import (
	"math/rand"
)

const (
	SwapBorderProbability = 0.4
)

func FormGroups(ratedParticipants []string, groupSize int) [][]string {
	return FormGroupsWithSwapProbability(ratedParticipants, groupSize, SwapBorderProbability)
}

func FormGroupsWithSwapProbability(ratedParticipants []string, groupSize int, swapBorderProbability float64) [][]string {
	n := len(ratedParticipants)
	if n == 0 {
		return [][]string{}
	}
	if groupSize <= 1 {
		groupSize = 2
	}
	swapBorderProbability = max(0, min(1, swapBorderProbability))

	participants := append([]string{}, ratedParticipants...)

	for i := 0; i < n-1; i++ {
		if rand.Float64() <= swapBorderProbability {
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

	lastIndex := len(result) - 1
	last := result[len(result)-1]
	if lastIndex >= 1 && len(last) < len(result[0])-1 {
		prev := result[lastIndex-1]
		tail := append(append([]string{}, prev...), result[lastIndex]...)

		newPrev := tail[:len(tail)/2+len(tail)%2]
		newLast := tail[len(tail)/2+len(tail)%2:]
		result[lastIndex-1] = newPrev
		result[lastIndex] = newLast
	}

	return result
}
