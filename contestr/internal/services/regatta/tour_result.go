package regatta

import (
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"slices"
	"strings"
)

const (
	OvertakePoints    = regatta.DefaultOvertakeBonus
	SolveInTimePoints = regatta.DefaultSolveInTimeBonus

	fullSolvePointsThreshold = 100
)

const (
	SubmissionStatusOK      string = "OK"
	SubmissionStatusPartial string = "PARTIAL"
)

const contestElapsedUnknown = -1

type Problem = int

type Participant = string

type SubmissionTime = int

type ScoredSubmission struct {
	Time   SubmissionTime
	Points int
	Status string
}

type ProblemState struct {
	BestPoints    int
	BestTime      SubmissionTime
	HasBest       bool
	FullSolveTime SubmissionTime
	HasFullSolve  bool
	Submissions   []ScoredSubmission
}

type ContestResult = map[Problem]ProblemState

type Group = []Participant

type TourResult struct {
	regatta.Tour
	SegmentStart          int
	ContestElapsedSeconds int
	ScoringSettings       regatta.ScoringSettings
	Results               map[Participant]ContestResult `json:"-"`
	ProblemsMapping       map[Problem]string            `json:"-"`
}

type ProblemResult struct {
	problemCode        string
	score              int
	lastSubmissionTime SubmissionTime
}

type ProblemCode = string

type ParticipantResult = map[ProblemCode]ProblemResult

type ResultsByParticipant = map[Participant]ParticipantResult

type ParticipantScoreOutcome struct {
	Results ParticipantResult
	Events  []regatta.RegattaEvent
}

type eventMeta struct {
	names map[string]string
}

func isSolvedInTime(segmentStart int, duration int, solveTime int) bool {
	return solveTime >= segmentStart && solveTime <= segmentStart+duration
}

func tourEnd(t *TourResult) int {
	return t.SegmentStart + t.DurationInSeconds
}

func tourEnded(t *TourResult) bool {
	return t.ContestElapsedSeconds == contestElapsedUnknown || t.ContestElapsedSeconds >= tourEnd(t)
}

func problemState(t *TourResult, participant Participant, problem Problem) (ProblemState, bool) {
	participantResults := t.Results[participant]
	if participantResults == nil {
		return ProblemState{}, false
	}
	state, ok := participantResults[problem]
	return state, ok
}

func pointsAtTime(state ProblemState, deadline int) int {
	points, _, _ := bestSubmissionAtTime(state, deadline)
	return points
}

func bestSubmissionAtTime(state ProblemState, deadline int) (int, SubmissionTime, bool) {
	bestPoints := 0
	bestTime := 0
	hasBest := false
	for _, sub := range state.Submissions {
		if sub.Time > deadline || sub.Points <= 0 {
			continue
		}
		if !hasBest || sub.Points > bestPoints || (sub.Points == bestPoints && sub.Time < bestTime) {
			bestPoints = sub.Points
			bestTime = sub.Time
			hasBest = true
		}
	}
	return bestPoints, bestTime, hasBest
}

func partialFullSolveWinner(t *TourResult, group Group, problem Problem) (Participant, bool, bool) {
	bestTime := 0
	var winners []Participant

	for _, participant := range group {
		state, ok := problemState(t, participant, problem)
		if !ok || !state.HasFullSolve || !isSolvedInTime(t.SegmentStart, t.DurationInSeconds, state.FullSolveTime) {
			continue
		}
		if len(winners) == 0 || state.FullSolveTime < bestTime {
			bestTime = state.FullSolveTime
			winners = []Participant{participant}
			continue
		}
		if state.FullSolveTime == bestTime {
			winners = append(winners, participant)
		}
	}

	if len(winners) == 1 {
		return winners[0], true, false
	}
	return "", false, len(winners) > 1
}

func partialTourEndWinner(t *TourResult, group Group, problem Problem) (Participant, bool) {
	if !tourEnded(t) {
		return "", false
	}

	end := tourEnd(t)
	bestPoints := 0
	var winners []Participant

	for _, participant := range group {
		state, ok := problemState(t, participant, problem)
		if !ok {
			continue
		}
		points := pointsAtTime(state, end)
		if points <= 0 {
			continue
		}
		if points > bestPoints {
			bestPoints = points
			winners = []Participant{participant}
			continue
		}
		if points == bestPoints {
			winners = append(winners, participant)
		}
	}

	if bestPoints > 0 && len(winners) == 1 {
		return winners[0], true
	}
	return "", false
}

func partialOvertakeWinner(t *TourResult, group Group, problem Problem) (Participant, bool, bool) {
	if len(group) <= 1 {
		return "", false, false
	}
	if winner, ok, tied := partialFullSolveWinner(t, group, problem); ok || tied {
		return winner, ok, false
	}
	if winner, ok := partialTourEndWinner(t, group, problem); ok {
		return winner, true, true
	}
	return "", false, false
}

func overtakeWinner(t *TourResult, participant Participant, problem Problem) bool {
	group := t.Groups[participant]
	winner, ok, _ := partialOvertakeWinner(t, group, problem)
	return ok && winner == participant
}

func scoreForProblem(t *TourResult, participant Participant, problem Problem) int {
	state, ok := problemState(t, participant, problem)
	if !ok || !state.HasBest {
		return 0
	}

	score := state.BestPoints
	if state.HasFullSolve &&
		isSolvedInTime(t.SegmentStart, t.DurationInSeconds, state.FullSolveTime) {
		score += t.ScoringSettings.SolveInTimeBonus
	}
	if overtakeWinner(t, participant, problem) {
		score += t.ScoringSettings.OvertakeBonus
	}
	return score
}

func solveEvent(
	meta eventMeta,
	t *TourResult,
	participant Participant,
	problem Problem,
	problemCode string,
	solveTime int,
) regatta.RegattaEvent {
	displayName := meta.names[participant]
	if displayName == "" {
		displayName = participant
	}

	inTime := isSolvedInTime(t.SegmentStart, t.DurationInSeconds, solveTime)
	firstInGroup := overtakeWinner(t, participant, problem)

	return regatta.RegattaEvent{
		Type:          regatta.EventTypeProblemSolved,
		TimeSec:       solveTime,
		ParticipantID: participant,
		DisplayName:   displayName,
		ProblemCode:   problemCode,
		TeamNumber:    t.GroupNumbers[participant],
		Points:        scoreForProblem(t, participant, problem),
		SolvedInTime:  inTime,
		FirstInGroup:  firstInGroup,
	}
}

func isFullSolveSubmission(sub ScoredSubmission, state ProblemState) bool {
	if !state.HasFullSolve {
		return false
	}
	return sub.Time == state.FullSolveTime &&
		sub.Points >= fullSolvePointsThreshold &&
		sub.Status == SubmissionStatusOK
}

func partialScoreEvent(
	meta eventMeta,
	t *TourResult,
	participant Participant,
	problem Problem,
	problemCode string,
	submissionTime SubmissionTime,
	points int,
) regatta.RegattaEvent {
	displayName := meta.names[participant]
	if displayName == "" {
		displayName = participant
	}

	return regatta.RegattaEvent{
		Type:          regatta.EventTypeProblemSolved,
		TimeSec:       submissionTime,
		ParticipantID: participant,
		DisplayName:   displayName,
		ProblemCode:   problemCode,
		TeamNumber:    t.GroupNumbers[participant],
		Points:        points,
		SolvedInTime:  isSolvedInTime(t.SegmentStart, t.DurationInSeconds, submissionTime),
	}
}

func rejectedEvent(
	meta eventMeta,
	t *TourResult,
	participant Participant,
	problem Problem,
	problemCode string,
	submissionTime SubmissionTime,
	status string,
) regatta.RegattaEvent {
	displayName := meta.names[participant]
	if displayName == "" {
		displayName = participant
	}

	return regatta.RegattaEvent{
		Type:          regatta.EventTypeProblemRejected,
		TimeSec:       submissionTime,
		ParticipantID: participant,
		DisplayName:   displayName,
		ProblemCode:   problemCode,
		TeamNumber:    t.GroupNumbers[participant],
		Points:        0,
		SolvedInTime:  isSolvedInTime(t.SegmentStart, t.DurationInSeconds, submissionTime),
		Verdict:       regatta.ShortSubmissionVerdict(status),
	}
}

func partialTourEndScoredEvent(
	meta eventMeta,
	t *TourResult,
	participant Participant,
	problem Problem,
	problemCode string,
	scoreTime int,
) regatta.RegattaEvent {
	displayName := meta.names[participant]
	if displayName == "" {
		displayName = participant
	}

	return regatta.RegattaEvent{
		Type:          regatta.EventTypeProblemOvertake,
		TimeSec:       scoreTime,
		ParticipantID: participant,
		DisplayName:   displayName,
		ProblemCode:   problemCode,
		TeamNumber:    t.GroupNumbers[participant],
		Points:        scoreForProblem(t, participant, problem),
		SolvedInTime:  isSolvedInTime(t.SegmentStart, t.DurationInSeconds, scoreTime),
	}
}

func tourEndPartialAwardKey(participant Participant, problemCode string) string {
	return participant + ":" + problemCode
}

func (t *TourResult) ParticipantScore(participant Participant) ParticipantScoreOutcome {
	logger.Infof(context.Background(), "Calculating participant score for participant %v in %v", participant, t.Name)

	participantResults := t.Results[participant]
	result := make(ParticipantResult, len(t.Problems))
	var events []regatta.RegattaEvent

	meta := eventMeta{names: nil}

	for _, problem := range t.Problems {
		problemCode := t.ProblemsMapping[problem]
		state, participantScored := participantResults[problem]
		if !participantScored {
			logger.Infof(context.Background(), "Participant %v did not score problem %v", participant, problem)
			result[problemCode] = ProblemResult{
				problemCode: problemCode,
				score:       0,
			}
			continue
		}

		if !state.HasBest {
			rejected := rejectedAttemptCount(state)
			logger.Infof(context.Background(), "Participant %v did not score problem %v", participant, problem)
			result[problemCode] = ProblemResult{
				problemCode: problemCode,
				score:       -rejected,
			}
			continue
		}

		score := scoreForProblem(t, participant, problem)
		logger.Infof(context.Background(),
			"Participant %v scored problem %v: score %v", participant, problem, score)

		result[problemCode] = ProblemResult{
			problemCode:        problemCode,
			score:              score,
			lastSubmissionTime: state.BestTime - t.SegmentStart,
		}

		if state.HasFullSolve {
			events = append(events, solveEvent(
				meta,
				t,
				participant,
				problem,
				problemCode,
				state.FullSolveTime,
			))
		}
	}

	return ParticipantScoreOutcome{Results: result, Events: events}
}

func (t *TourResult) Export() ResultsByParticipant {
	result := make(ResultsByParticipant)

	for participant := range t.Groups {
		result[participant] = t.ParticipantScore(participant).Results
	}

	return result
}

type Run struct {
	UserID string
	ProbID int
	Time   int
	Status string
	Points int
}

func CalculateResult(tc regatta.Tour, segmentStart int, runs []Run) *TourResult {
	return CalculateResultWithSettings(tc, segmentStart, runs, regatta.DefaultScoringSettings())
}

func CalculateResultWithSettings(tc regatta.Tour, segmentStart int, runs []Run, settings regatta.ScoringSettings) *TourResult {
	return CalculateResultWithSettingsAt(tc, segmentStart, runs, settings, contestElapsedUnknown)
}

func CalculateResultWithSettingsAt(
	tc regatta.Tour,
	segmentStart int,
	runs []Run,
	settings regatta.ScoringSettings,
	contestElapsedSeconds int,
) *TourResult {
	settings = regatta.NormalizeScoringSettings(settings)
	return &TourResult{
		Tour:                  tc,
		SegmentStart:          segmentStart,
		ContestElapsedSeconds: contestElapsedSeconds,
		ScoringSettings:       settings,
		Results:               calcSubmissions(runs, settings),
		ProblemsMapping:       tc.ProblemsIDsToNameMapping(tc.Problems),
	}
}

func calcSubmissions(submissions []Run, settings regatta.ScoringSettings) map[Participant]ContestResult {
	results := make(map[Participant]ContestResult)

	for _, submission := range submissions {
		points, keep := normalizedRunPoints(submission, settings)
		if !keep {
			continue
		}

		logger.Infof(context.Background(), "Scored submission! UserID: %v, ProbID: %v, Time: %v, Points: %v",
			submission.UserID, submission.ProbID, submission.Time, points)

		participant := submission.UserID
		problem := submission.ProbID

		participantResults, ok := results[participant]
		if !ok {
			participantResults = make(ContestResult)
			results[participant] = participantResults
		}

		state := participantResults[problem]
		scored := ScoredSubmission{
			Time:   submission.Time,
			Points: points,
			Status: submission.Status,
		}
		state.Submissions = append(state.Submissions, scored)
		if points > 0 && (!state.HasBest ||
			points > state.BestPoints ||
			(points == state.BestPoints && submission.Time < state.BestTime)) {
			state.HasBest = true
			state.BestPoints = points
			state.BestTime = submission.Time
		}
		if points >= fullSolvePointsThreshold && (!state.HasFullSolve || submission.Time < state.FullSolveTime) {
			state.HasFullSolve = true
			state.FullSolveTime = submission.Time
		}
		participantResults[problem] = state
	}

	return results
}

func isRejectedSubmission(sub ScoredSubmission) bool {
	if regatta.IsIgnorableSubmissionStatus(sub.Status) {
		return false
	}
	if sub.Points > 0 {
		return false
	}
	switch sub.Status {
	case SubmissionStatusOK, SubmissionStatusPartial:
		return false
	default:
		return true
	}
}

func rejectedAttemptCount(state ProblemState) int {
	count := 0
	for _, sub := range state.Submissions {
		if isRejectedSubmission(sub) {
			count++
		}
	}
	return count
}

func normalizedRunPoints(run Run, _ regatta.ScoringSettings) (int, bool) {
	if run.Status == SubmissionStatusOK {
		return 100, true
	}

	points := run.Points
	if points < 0 {
		points = 0
	}
	if points > 0 || run.Status == SubmissionStatusPartial {
		return points, true
	}
	return 0, isRejectedSubmission(ScoredSubmission{Status: run.Status})
}

func tourForProblem(tours []regatta.Tour, probID Problem) *regatta.Tour {
	for i := range tours {
		tour := &tours[i]
		if tour.IsPause {
			continue
		}
		if slices.Contains(tour.Problems, probID) {
			return tour
		}
	}
	return nil
}

func BuildContestEvents(
	tours []regatta.Tour,
	runs []Run,
	names map[string]string,
	settingsArg ...regatta.ScoringSettings,
) []regatta.RegattaEvent {
	return BuildContestEventsAt(tours, runs, names, contestElapsedUnknown, settingsArg...)
}

func BuildContestEventsAt(
	tours []regatta.Tour,
	runs []Run,
	names map[string]string,
	contestElapsedSeconds int,
	settingsArg ...regatta.ScoringSettings,
) []regatta.RegattaEvent {
	if len(tours) == 0 {
		return nil
	}

	settings := regatta.DefaultScoringSettings()
	if len(settingsArg) > 0 {
		settings = regatta.NormalizeScoringSettings(settingsArg[0])
	}

	sorted := regatta.SortToursBySequence(tours)
	offsets := regatta.SegmentOffsets(sorted)
	meta := eventMeta{names: names}
	var allEvents []regatta.RegattaEvent

	for _, tour := range sorted {
		if tour.IsPause {
			continue
		}
		segmentStart := offsets[tour.Sequence].Start
		tr := CalculateResultWithSettingsAt(tour, segmentStart, runs, settings, contestElapsedSeconds)
		tourEndEvents := buildPartialOvertakeEvents(meta, tr)
		tourEndAwards := make(map[string]int, len(tourEndEvents))
		for _, e := range tourEndEvents {
			tourEndAwards[tourEndPartialAwardKey(e.ParticipantID, e.ProblemCode)] = e.TimeSec
		}

		for participant := range tour.Groups {
			for _, problem := range tour.Problems {
				problemCode := tr.ProblemsMapping[problem]
				if problemCode == "" {
					continue
				}

				state, ok := problemState(tr, participant, problem)
				if ok {
					for _, sub := range state.Submissions {
						if isRejectedSubmission(sub) {
							allEvents = append(allEvents,
								rejectedEvent(meta, tr, participant, problem, problemCode, sub.Time, sub.Status))
							continue
						}
						if sub.Points > 0 && !isFullSolveSubmission(sub, state) {
							key := tourEndPartialAwardKey(participant, problemCode)
							if awardTime, skip := tourEndAwards[key]; skip && sub.Time == awardTime {
								continue
							}
							allEvents = append(allEvents,
								partialScoreEvent(meta, tr, participant, problem, problemCode, sub.Time, sub.Points))
						}
					}
					if state.HasFullSolve {
						allEvents = append(allEvents,
							solveEvent(meta, tr, participant, problem, problemCode, state.FullSolveTime))
					}
				}
			}
		}

		allEvents = append(allEvents, tourEndEvents...)
	}

	slices.SortFunc(allEvents, func(a, b regatta.RegattaEvent) int {
		if a.TimeSec != b.TimeSec {
			return a.TimeSec - b.TimeSec
		}
		if a.ParticipantID != b.ParticipantID {
			return strings.Compare(a.ParticipantID, b.ParticipantID)
		}
		return strings.Compare(a.ProblemCode, b.ProblemCode)
	})

	return allEvents
}

func buildPartialOvertakeEvents(meta eventMeta, t *TourResult) []regatta.RegattaEvent {
	var events []regatta.RegattaEvent
	seen := make(map[string]bool)

	for participant, group := range t.Groups {
		groupNumber := t.GroupNumbers[participant]
		for _, problem := range t.Problems {
			key := fmt.Sprintf("%d:%d", groupNumber, problem)
			if seen[key] {
				continue
			}
			seen[key] = true

			winner, ok, tourEndFallback := partialOvertakeWinner(t, group, problem)
			if !ok || !tourEndFallback {
				continue
			}

			problemCode := t.ProblemsMapping[problem]
			if problemCode == "" {
				continue
			}
			state, ok := problemState(t, winner, problem)
			if !ok {
				continue
			}
			if state.HasFullSolve {
				continue
			}
			_, scoreTime, ok := bestSubmissionAtTime(state, tourEnd(t))
			if !ok {
				continue
			}
			events = append(events, partialTourEndScoredEvent(meta, t, winner, problem, problemCode, scoreTime))
		}
	}

	return events
}
