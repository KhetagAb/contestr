package regatta

import "time"

type Contest struct {
	ContestID   int       `bson:"contest_id"`
	ContestName string    `bson:"contest_name"`
	System      string    `bson:"system"`
	StartTime   time.Time `bson:"start_time"`
	LastUpdated time.Time `bson:"last_updated"`

	Participants []ContestParticipant `bson:"participants"`
	Submissions  []ContestSubmission  `bson:"submissions"`
}

type ContestParticipant struct {
	ID          string `bson:"id"`
	DisplayName string `bson:"display_name"`
	OriginalID  string `bson:"original_id"`
}

type ContestSubmission struct {
	ParticipantID     string `bson:"participant_id"`
	ProblemID         int    `bson:"problem_id"`
	Time              int    `bson:"time"`
	Status            string `bson:"status"`
	OriginalProblemID string `bson:"original_problem_id"`
}
