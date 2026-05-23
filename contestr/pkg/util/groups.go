package util

import (
	"cmp"
	"math/rand"
	"slices"
)

func FormGroupsWithSwapProbability(ratedParticipants []string, groupSize int, shuffleProbability float64) [][]string {
	n := len(ratedParticipants)
	if n == 0 {
		return [][]string{}
	}
	if groupSize <= 1 {
		groupSize = 2
	}
	shuffleProbability = max(0, min(1, shuffleProbability))

	participants := shuffleParticipantsWithRNG(ratedParticipants, shuffleProbability, nil)

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

func shuffleParticipantsWithRNG(ratedParticipants []string, p float64, rng *rand.Rand) []string {
	n := len(ratedParticipants)
	if p <= 0 {
		return append([]string{}, ratedParticipants...)
	}

	randomFloat := rand.Float64
	shuffle := rand.Shuffle
	if rng != nil {
		randomFloat = rng.Float64
		shuffle = rng.Shuffle
	}

	participants := append([]string{}, ratedParticipants...)
	if p >= 1 {
		shuffle(len(participants), func(i, j int) {
			participants[i], participants[j] = participants[j], participants[i]
		})
		return participants
	}

	type keyed struct {
		participant string
		key         float64
	}
	rankScale := 1.0
	if n > 1 {
		rankScale = float64(n-1) / (1 - p)
	}
	keys := make([]keyed, n)
	for i, participant := range ratedParticipants {
		keys[i] = keyed{
			participant: participant,
			key:         float64(i) + p*randomFloat()*rankScale,
		}
	}
	slices.SortFunc(keys, func(a, b keyed) int {
		return cmp.Compare(a.key, b.key)
	})
	for i := range keys {
		participants[i] = keys[i].participant
	}
	return participants
}
