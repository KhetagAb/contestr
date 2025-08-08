package util

import "math/rand"

const (
	SWAP_BORDER_PROBABILITY = 0.4
)

func FormGroups(ratedParticipants []int, groupSize int) [][]int {
	n := len(ratedParticipants)
	result := make([][]int, 0, n/groupSize+1)

	for i := 0; i < n-1; i++ {
		if rand.Float64() <= SWAP_BORDER_PROBABILITY {
			ratedParticipants[i], ratedParticipants[i+1] = ratedParticipants[i+1], ratedParticipants[i]
		}
	}

	for i := 0; i < n; i += groupSize {
		intervalEnd := i + groupSize
		if intervalEnd > n {
			intervalEnd = n
		}

		if intervalEnd == n-1 {
			result = append(result, ratedParticipants[i:intervalEnd-1])
			result = append(result, ratedParticipants[intervalEnd-1:n])
		} else {
			result = append(result, ratedParticipants[i:intervalEnd])
		}
	}

	return result
}
