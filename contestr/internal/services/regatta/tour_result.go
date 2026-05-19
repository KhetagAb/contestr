package regatta

import (
	"contestr/pkg/logger"
	"contestr/pkg/regatta"
	"context"
	"slices"
)

const (
	OvertakePoints    = 5
	SolvePoints       = 5
	SolveInTimePoints = 5
)

const (
	SubmissionStatusOK string = "OK"
)

type Problem = int

type Participant = string

type SubmissionTime = int

type ContestResult = map[Problem]SubmissionTime

type Group = []Participant

type TourResult struct {
	regatta.Tour
	SegmentStart    int
	Results         map[Participant]ContestResult `json:"-"`
	ProblemsMapping map[Problem]string            `json:"-"`
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
	return solveTime <= segmentStart+duration
}

func overtakeCount(t *TourResult, participant Participant, problem Problem, solveTime int) int {
	group := t.Groups[participant]
	count := 0
	for _, opponent := range group {
		if participant == opponent {
			continue
		}
		opponentResults := t.Results[opponent]
		opponentSolveTime, opponentSolved := opponentResults[problem]
		if !opponentSolved || solveTime < opponentSolveTime {
			count++
		}
	}
	return count
}

func isFirstInGroupRetrospective(t *TourResult, participant Participant, problem Problem, solveTime int) bool {
	group := t.Groups[participant]
	if len(group) <= 1 {
		return false
	}
	return overtakeCount(t, participant, problem, solveTime) == len(group)-1
}

func scoreForSolve(t *TourResult, participant Participant, problem Problem, solveTime int) int {
	score := SolvePoints
	if isSolvedInTime(t.SegmentStart, t.DurationInSeconds, solveTime) {
		score += SolveInTimePoints
	}
	if isFirstInGroupRetrospective(t, participant, problem, solveTime) {
		score += OvertakePoints
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
	firstInGroup := isFirstInGroupRetrospective(t, participant, problem, solveTime)

	return regatta.RegattaEvent{
		Type:          regatta.EventTypeProblemSolved,
		TimeSec:       solveTime,
		ParticipantID: participant,
		DisplayName:   displayName,
		ProblemCode:   problemCode,
		TeamNumber:    t.GroupNumbers[participant],
		Points:        scoreForSolve(t, participant, problem, solveTime),
		SolvedInTime:  inTime,
		FirstInGroup:  firstInGroup,
	}
}

func (t *TourResult) ParticipantScore(participant Participant) ParticipantScoreOutcome {
	logger.Infof(context.Background(), "Calculating participant score for participant %v in %v", participant, t.Name)

	participantResults := t.Results[participant]
	result := make(ParticipantResult, len(t.Problems))
	var events []regatta.RegattaEvent

	meta := eventMeta{names: nil}

	for _, problem := range t.Problems {
		problemCode := t.ProblemsMapping[problem]

		participantSolveTime, participantSolved := participantResults[problem]
		if !participantSolved {
			logger.Infof(context.Background(), "Participant %v did not solve problem %v", participant, problem)
			result[problemCode] = ProblemResult{
				problemCode: problemCode,
				score:       0,
			}
			continue
		}

		score := scoreForSolve(t, participant, problem, participantSolveTime)
		logger.Infof(context.Background(),
			"Participant %v solved problem %v: score %v", participant, problem, score)

		result[problemCode] = ProblemResult{
			problemCode:        problemCode,
			score:              score,
			lastSubmissionTime: participantSolveTime - t.SegmentStart,
		}

		events = append(events, solveEvent(
			meta,
			t,
			participant,
			problem,
			problemCode,
			participantSolveTime,
		))
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
}

func CalculateResult(tc regatta.Tour, segmentStart int, runs []Run) *TourResult {
	return &TourResult{
		Tour:            tc,
		SegmentStart:    segmentStart,
		Results:         calcSubmissions(runs),
		ProblemsMapping: tc.ProblemsIDsToNameMapping(tc.Problems),
	}
}

func calcSubmissions(submissions []Run) map[Participant]ContestResult {
	results := make(map[Participant]ContestResult)

	for _, submission := range submissions {
		if submission.Status != SubmissionStatusOK {
			continue
		}

		logger.Infof(context.Background(), "OK! UserID: %v, ProbID: %v, Time: %v", submission.UserID, submission.ProbID, submission.Time)

		participant := submission.UserID
		problem := submission.ProbID

		participantResults, ok := results[participant]
		if !ok {
			results[participant] = map[Problem]int{
				problem: submission.Time,
			}
		} else {
			if _, alreadySolved := participantResults[problem]; !alreadySolved {
				participantResults[problem] = submission.Time
			}
		}
	}

	return results
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

func cloneResults(src map[Participant]ContestResult) map[Participant]ContestResult {
	dst := make(map[Participant]ContestResult, len(src))
	for participant, problems := range src {
		cp := make(ContestResult, len(problems))
		for problem, t := range problems {
			cp[problem] = t
		}
		dst[participant] = cp
	}
	return dst
}

func applySolve(results map[Participant]ContestResult, participant Participant, problem Problem, solveTime int) {
	participantResults, ok := results[participant]
	if !ok {
		results[participant] = ContestResult{problem: solveTime}
		return
	}
	if _, alreadySolved := participantResults[problem]; !alreadySolved {
		participantResults[problem] = solveTime
	}
}

func isFirstAC(results map[Participant]ContestResult, participant Participant, problem Problem) bool {
	participantResults, ok := results[participant]
	if !ok {
		return true
	}
	_, solved := participantResults[problem]
	return !solved
}

func BuildContestEvents(
	tours []regatta.Tour,
	runs []Run,
	names map[string]string,
) []regatta.RegattaEvent {
	if len(tours) == 0 {
		return nil
	}

	sorted := regatta.SortToursBySequence(tours)
	offsets := regatta.SegmentOffsets(sorted)

	okRuns := make([]Run, 0, len(runs))
	for _, run := range runs {
		if run.Status == SubmissionStatusOK {
			okRuns = append(okRuns, run)
		}
	}
	slices.SortFunc(okRuns, func(a, b Run) int {
		if a.Time != b.Time {
			return a.Time - b.Time
		}
		if a.UserID != b.UserID {
			if a.UserID < b.UserID {
				return -1
			}
			return 1
		}
		return a.ProbID - b.ProbID
	})

	meta := eventMeta{names: names}
	results := make(map[Participant]ContestResult)
	var allEvents []regatta.RegattaEvent

	for _, run := range okRuns {
		if !isFirstAC(results, run.UserID, run.ProbID) {
			continue
		}

		tour := tourForProblem(sorted, run.ProbID)
		if tour == nil {
			applySolve(results, run.UserID, run.ProbID, run.Time)
			continue
		}

		snapshot := cloneResults(results)
		segmentStart := offsets[tour.Sequence].Start
		tr := &TourResult{
			Tour:            *tour,
			SegmentStart:    segmentStart,
			Results:         snapshot,
			ProblemsMapping: tour.ProblemsIDsToNameMapping(tour.Problems),
		}

		problemCode, ok := tr.ProblemsMapping[run.ProbID]
		if !ok || problemCode == "" {
			applySolve(results, run.UserID, run.ProbID, run.Time)
			continue
		}

		allEvents = append(allEvents, solveEvent(meta, tr, run.UserID, run.ProbID, problemCode, run.Time))

		applySolve(results, run.UserID, run.ProbID, run.Time)
	}

	return allEvents
}
