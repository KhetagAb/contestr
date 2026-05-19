package regatta

import "time"

const (
	TourStatusStarted      = "started"
	TourStatusNext         = "next"
	TourStatusStarting = "starting"
	TourStatusPlanned  = "planned"
)

type TourMeta struct {
	TourNumber int    `json:"tour_number"`
	Status     string `json:"status"`
}

type SaveTimetableRequest struct {
	TourDurations    []int `json:"tour_durations"`
	AutoStartEnabled bool  `json:"auto_start_enabled"`
}

type TimetableView struct {
	ContestId           int          `json:"contest_id"`
	ContestStartTime    *time.Time   `json:"contest_start_time,omitempty"`
	ServerNow           time.Time    `json:"server_now"`
	ElapsedSeconds      int          `json:"elapsed_seconds"`
	TourTimes           []TourConfig `json:"tour_times"`
	ToursMeta           []TourMeta   `json:"tours_meta"`
	NextTourNumber      *int         `json:"next_tour_number"`
	AutoStartEnabled   bool `json:"auto_start_enabled"`
	AutoStartAvailable bool `json:"auto_start_available"`
}
