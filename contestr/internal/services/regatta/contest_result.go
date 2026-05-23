package regatta

import (
	"contestr/pkg/regatta"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
)

// TODO: при введении снепшотов результатов пересчитывать только новые туры,
// накладывая их поверх последнего снепшота вместо полного пересчёта с нуля.
func (s *Regatta) GetContestResult(ctx context.Context, contestID int) (regatta.ContestStandings, error) {
	contest, err := s.contestRepo.GetByContestID(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to get contest %d: %w", contestID, err)
	}

	participantsMap, err := s.contestRepo.GetParticipants(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to get participants: %w", err)
	}

	tours, err := s.loadToursSorted(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, err
	}

	currentTime := time.Now()
	currentElapsedSeconds := max(0, int(time.Since(contest.StartTime).Seconds()))
	standings := regatta.ContestStandings{
		ContestId:        contestID,
		ContestName:      contest.ContestName,
		CurrentTime:      currentTime,
		ContestStartTime: contest.StartTime,
		Rows:             []regatta.ContestRow{},
		Events:           []regatta.RegattaEvent{},
	}

	if len(tours) == 0 {
		contestRows := []regatta.ContestRow{}
		for _, participant := range contest.Participants {
			contestRows = append(contestRows, regatta.ContestRow{
				DisplayName:    participant.DisplayName,
				UserID:         participant.ID,
				ProblemResults: []regatta.ProblemResult{},
				SolvedProblems: 0,
				TeamNumber:     0,
				TotalScore:     0,
			})
		}
		standings.Rows = contestRows
		standings.CurrentTourStartTime = 0
		standings.CurrentTourDuration = 0
		return standings, nil
	}

	submissions, err := s.contestRepo.GetSubmissions(ctx, contestID)
	if err != nil {
		return regatta.ContestStandings{}, fmt.Errorf("failed to get submissions: %w", err)
	}

	runs := convertSubmissionsToRuns(submissions)
	scoringSettings := regatta.NormalizeScoringSettings(contest.ScoringSettings)

	var contestRows []regatta.ContestRow
	contestStandingsByParticipants := make(ResultsByParticipant)
	participantTotal := make(map[regatta.Participant]int)

	offsets := regatta.SegmentOffsets(tours)

	for _, tour := range tours {
		if tour.IsPause {
			continue
		}
		segmentStart := offsets[tour.Sequence].Start
		result := CalculateResultWithSettingsAt(tour, segmentStart, runs, scoringSettings, currentElapsedSeconds).Export()

		for participant, participantResult := range result {
			_, was := contestStandingsByParticipants[participant]
			if !was {
				contestStandingsByParticipants[participant] = make(ParticipantResult)
			}

			for problem, problemResult := range participantResult {
				contestStandingsByParticipants[participant][problem] = problemResult
				if problemResult.score > 0 {
					participantTotal[participant] += problemResult.score
				}
			}
		}
	}

	for participant, participantResult := range contestStandingsByParticipants {
		displayName := participantsMap[participant]
		if displayName == "" {
			displayName = participant
		}

		teamNumber := 0
		if teamTour := competitiveTourForTeamDisplay(tours, currentElapsedSeconds); teamTour != nil {
			teamNumber = teamTour.GroupNumbers[participant]
		}

		contestRows = append(contestRows, regatta.ContestRow{
			DisplayName:    displayName,
			ProblemResults: getProblemResults(participantResult),
			SolvedProblems: getSolvedProblemsCount(participantResult),
			TeamNumber:     teamNumber,
			TotalScore:     participantTotal[participant],
			UserID:         participant,
		})
	}

	slices.SortFunc(contestRows, func(row1, row2 regatta.ContestRow) int {
		if row1.TotalScore != row2.TotalScore {
			return row2.TotalScore - row1.TotalScore
		}
		return strings.Compare(row1.DisplayName, row2.DisplayName)
	})
	standings.Rows = contestRows

	if len(tours) > 0 {
		lastTour := tours[len(tours)-1]
		lastOffset := offsets[lastTour.Sequence]
		standings.CurrentTourStartTime = int(contest.StartTime.Unix()) + lastOffset.Start
		standings.CurrentTourDuration = lastTour.DurationInSeconds / 60
		standings.IsPauseBreak = lastTour.IsPause
		standings.Events = BuildContestEventsAt(tours, runs, participantsMap, currentElapsedSeconds, scoringSettings)
	}

	return standings, nil
}

func convertSubmissionsToRuns(submissions []regatta.ContestSubmission) []Run {
	runs := make([]Run, 0, len(submissions))
	for _, sub := range submissions {
		runs = append(runs, Run{
			UserID: sub.ParticipantID,
			ProbID: sub.ProblemID,
			Time:   sub.Time,
			Status: sub.Status,
			Points: sub.Points,
		})
	}
	return runs
}

func getProblemResults(result ParticipantResult) []regatta.ProblemResult {
	var results []regatta.ProblemResult
	for _, problemResult := range result {
		results = append(results, regatta.ProblemResult{
			ProblemCode:        problemResult.problemCode,
			Score:              problemResult.score,
			LastSubmissionTime: problemResult.lastSubmissionTime,
		})
	}
	slices.SortFunc(results, func(result1, result2 regatta.ProblemResult) int {
		return strings.Compare(result1.ProblemCode, result2.ProblemCode)
	})
	return results
}

// competitiveTourForTeamDisplay — тур, по которому в таблице показывается группа:
// идущий сейчас или последний уже начатый (на паузе между турами).
func competitiveTourForTeamDisplay(tours []regatta.Tour, elapsed int) *regatta.Tour {
	if len(tours) == 0 {
		return nil
	}
	sorted := regatta.SortToursBySequence(tours)
	offsets := regatta.SegmentOffsets(tours)
	if seq, ok := regatta.ActiveSequence(tours, elapsed); ok {
		for i := range sorted {
			if sorted[i].Sequence == seq && !sorted[i].IsPause {
				return &sorted[i]
			}
		}
	}
	var lastStarted *regatta.Tour
	for i := range sorted {
		t := &sorted[i]
		if t.IsPause {
			continue
		}
		if elapsed >= offsets[t.Sequence].End {
			lastStarted = t
		}
	}
	return lastStarted
}

func getSolvedProblemsCount(result ParticipantResult) int {
	count := 0
	for _, problemResult := range result {
		if problemResult.score > 0 {
			count++
		}
	}
	return count
}
